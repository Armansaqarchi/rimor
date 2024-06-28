package mapreducer

import (
	"rimor/pkg/engine/dictionary/record"
	"rimor/pkg/engine/dictionary/xindex"
	"rimor/pkg/engine/inverter/mapreducer/segment"
	"sort"
)

type Inverter struct {
	Out *xindex.Xindex
}

func NewInverter() Inverter {
	return Inverter{
		Out: &xindex.Xindex{
			Records:     make([]record.Recorder, 0),
			ChampionNum: 1024,
		},
	}
}

func (inv *Inverter) Serve(Input *segment.Segment) {
	combinedFragment := segment.Fragment{
		Pairs: make([]*segment.Pair, 0),
	}
	for _, frag := range Input.Fragments {
		combinedFragment.Pairs = append(combinedFragment.Pairs, frag.Pairs...)
	}

	if len(combinedFragment.Pairs) == 0 {
		return
	}

	sort.Sort(combinedFragment)
	inv.Out.Records = append(inv.Out.Records, record.NewRecord(combinedFragment.Pairs[0].Term, combinedFragment.Pairs[0].Doc))
	currentRec := inv.Out.Records[0]
	for _, t := range combinedFragment.Pairs[1:] {
		if t.Term == currentRec.GetTerm() {
			if currentRec.GetLast().GetDocID() == t.Doc {
				currentRec.GetLast().IncreaseTF()
			} else {
				currentRec.AddToPosting(record.NewPostingListElem(t.Doc, nil))
				currentRec.IncreaseDF()
			}
			continue
		}
		currentRec = record.NewRecord(
			t.Term,
			t.Doc,
		)
		inv.Out.Records = append(inv.Out.Records, currentRec)
	}
}
