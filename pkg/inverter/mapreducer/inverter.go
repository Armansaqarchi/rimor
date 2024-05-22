package mapreducer

import (
	"rimor/pkg/dictionary/record"
	"rimor/pkg/dictionary/xindex"
	"rimor/pkg/inverter/mapreducer/segment"
	"sort"
)

type Inverter struct {
	Input *segment.Segment
	Out xindex.Xindex
}


func NewInverter(input *segment.Segment) *Inverter{
	return &Inverter{
		Input: input,
		Out: xindex.Xindex{
			Records: make([]record.Recorder, 0),
		},
	}
}


func (inv *Inverter) Serve() {
	combinedFragment := segment.Fragment{
		Pairs: make([]*segment.Pair, 0),
	}
	for _, frag := range inv.Input.Fragments {
		combinedFragment.Pairs = append(combinedFragment.Pairs, frag.Pairs...)
	}
	sort.Sort(combinedFragment)
	inv.Out.Records = append(inv.Out.Records, record.NewRecord(
		combinedFragment.Pairs[0].Term,
	))
	currentRec := inv.Out.Records[0]
	for _, t := range combinedFragment.Pairs[1:] {
		if t.Term == currentRec.GetTerm() {
			if currentRec.GetLast().GetDocID() == t.Doc{
				currentRec.GetLast().IncreaseTF()
			} else{
				currentRec.AddToPosting(record.NewPostingListElem(t.Doc, nil))
			}
			continue
		}
		currentRec = record.NewRecord(
			t.Term,
		)
		inv.Out.Records = append(inv.Out.Records, currentRec)
	}
}