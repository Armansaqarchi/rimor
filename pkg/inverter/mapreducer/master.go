package mapreducer

import (
	record "rimor/pkg/dictionary/record"
	index "rimor/pkg/dictionary/xindex"
	segment "rimor/pkg/inverter/mapreducer/segment"
	preprocessing "rimor/pkg/preprocessing"
	"sync"
)



type Master struct {
	splits int
	parsersNum int
	invertersNum int
	Segments []*segment.Segment
	wg *sync.WaitGroup
	inverters []Inverter
}


func NewMaster(invertersNum int, splits int, parsersNum int) *Master{
	m := Master{
		splits: splits,
		invertersNum: invertersNum,
		parsersNum: parsersNum,
	}
	for i := 0; i< parsersNum; i++ {
		s := segment.Segment{}
		for i := 0; i < invertersNum; i++{
			f := segment.Fragment{}
			f.Pairs = make([]*segment.Pair, 0)
			s.Fragments = append(s.Fragments, &f)
		}
		m.Segments = append(m.Segments, &s)
	}
	m.wg = &sync.WaitGroup{}
	for i := 0; i < m.invertersNum; i++ {
		m.inverters = append(m.inverters, NewInverter())
	}
	return &m
}





func (m *Master) MapReduce(tk preprocessing.TkDocumentCollection) index.Xindex{

	var splits chan *preprocessing.TkDocumentCollection = make(chan *preprocessing.TkDocumentCollection, 10)
	for i := 0; i < m.parsersNum; i++ {
		m.wg.Add(1)
		go func(splits chan *preprocessing.TkDocumentCollection, order int) {
			p := NewParser(
				m.invertersNum,
				m.Segments[order],
			)
			defer m.wg.Done()
			for split := range splits {
				p.Serve(split)
			}
		}(splits, i)
	}
	documentsLen := len(tk.DocList)
	splitLen :=  documentsLen / m.splits
	for i := 0; i < m.splits; i++ {
		split := tk.DocList[i * splitLen: (i + 1) * splitLen]
		splits <- &preprocessing.TkDocumentCollection{
			DocList: split,
		}
	}

	if splitLen * m.splits < documentsLen{
		split := tk.DocList[splitLen * m.splits: documentsLen]
		splits <- &preprocessing.TkDocumentCollection{
			DocList: split,
		}
	}

	close(splits)
	m.wg.Wait()
	for i := 0; i < len(m.inverters); i++ {
		m.wg.Add(1)
		go func (order int) {
			defer m.wg.Done()
			m.inverters[order].Serve(m.GetVerticalSegment(order))
		}(i)
	}

	m.wg.Wait()


	x := index.Xindex{
		Records: make([]record.Recorder, 0),
		DocNum: len(tk.DocList),
	}

	for i := 0; i < len(m.inverters); i++ {
		x.Records = append(x.Records, m.inverters[i].Out.Records...) 
	}
	

	return x
}



/* suppose that segments for the array are defined as follows:
array = [ [frag11] [frag12] [frag13] frag[14] frag[15] ]
		| [frag21] [frag22] [frag23] frag[24] frag[25] |
		| [frag31] [frag32] [frag33] frag[34] frag[35] |
		[ [frag41] [frag42] [frag43] frag[44] frag[45] ]


now, calling getVerticalSegment returns an 1d array of fragments of column "num"


example: GetVerticalSegment(2) would return [[frag12] [frag22] [frag32] [frag42]]

this actually maintains hierarchial implementation so that inverter workers could have a single 
Segment to parse instead of list of segments that not all of them are used

*/
func (m *Master) GetVerticalSegment(num int) *segment.Segment{
	if num >= len(m.Segments[0].Fragments) {
		return nil
	}
	segment := segment.Segment{}
	for i := 0; i < len(m.Segments); i++ {
		segment.Fragments = append(segment.Fragments, m.Segments[i].Fragments[num]) 
	}
	return &segment
}