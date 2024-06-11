package mapreducer

import (
	"sync"
	preprocessing "rimor/pkg/preprocessing"
	segment "rimor/pkg/inverter/mapreducer/segment"
)



type Master struct {
	Splits int
	Segments []*segment.Segment
	wg *sync.WaitGroup
	parsers int
	inverters int
}




func (m *Master) MapReduce(tk preprocessing.TkDocumentCollection) {

	var splits chan *preprocessing.TkDocumentCollection

	for i := 0; i < m.Splits; i++ {
		m.wg.Add(1)
		p := NewParser(
			m.inverters,
			m.Segments[i],
		)
		go func(splits chan *preprocessing.TkDocumentCollection) {
			for split := range splits {
				p.Serve(split)
			}
		}(splits)
	}
	documentsLen := len(tk.DocList)
	splitLen :=  documentsLen / m.Splits
	for i := 0; i < m.Splits; i++ {
		split := tk.DocList[i * splitLen: (i + 1) * splitLen]
		splits <- &preprocessing.TkDocumentCollection{
			DocList: split,
		}
	}
	if splitLen * m.Splits < documentsLen{
		split := tk.DocList[splitLen * m.Splits: documentsLen]
		splits <- &preprocessing.TkDocumentCollection{
			DocList: split,
		}
	}


	


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
func (m *Master) GetVerticalSegment(num int) *Segment{
	if num >= len(sm.Segments[0].Fragments) {
		return nil
	}
	segment := Segment{}
	for i := 0; i < len(sm.Segments); i++ {
		segment.Fragments = append(segment.Fragments, sm.Segments[i].Fragments[num]) 
	}
	return &segment
}