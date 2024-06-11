package mapreducer

import (
	consts "rimor/pkg/consts"
	segment "rimor/pkg/inverter/mapreducer/segment"
	preprocessing "rimor/pkg/preprocessing"
)

type Parser struct {
	runeGroups []rune
	Out *segment.Segment
}


func NewParser(groups int, Out *segment.Segment) Parser{
	criterion := consts.Criterion
	length := len(criterion)
	step :=  length / groups

	parser := Parser{
		runeGroups: make([]rune, 0),
		Out: Out,
	}	
	for i := 1; i <= groups + 1; i++{
		parser.runeGroups = append(parser.runeGroups, consts.Criterion[min(groups * step, length) - 1])
	}
}

func(p *Parser) Serve(Input *preprocessing.TkDocumentCollection) {
	for _, d := range Input.DocList{
		for _, t := range d.TokenzedDocContent {
			p.AddTokenToFragment(t, int64(d.Id))
		}
	}
}


func (p *Parser) AddTokenToFragment(t string, d int64) {
	if len(t) > 0 {
		// log if you want
	}
	for idx, r := range p.runeGroups {
		if rune(t[0]) <= r {
			p.Out.Fragments[idx].AddPair(&segment.Pair{
				Term: t,
				Doc: d,
			})
		}
	}
}