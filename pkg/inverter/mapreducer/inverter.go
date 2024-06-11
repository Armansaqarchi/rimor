package mapreducer

import (
	"rimor/pkg/dictionary/record"
	"rimor/pkg/dictionary/xindex"
	"rimor/pkg/inverter/mapreducer/segment"
	"sort"
)

type Inverter struct {
	Out *xindex.Xindex
}

func NewInverter(out *xindex.Xindex) *Inverter{
	return &Inverter{
		Out: out,
	}
}


func (inv *Inverter) Serve(Input *segment.Segment) {
	combinedFragment := segment.Fragment{
		Pairs: make([]*segment.Pair, 0),
	}
	for _, frag := range Input.Fragments {
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