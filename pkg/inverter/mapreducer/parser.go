package mapreducer

import (
	consts "rimor/pkg/consts"
	segment "rimor/pkg/inverter/mapreducer/segment"
	preprocessing "rimor/pkg/preprocessing"
)

type Parser struct {
	runeGroups []rune
	Input *preprocessing.TkDocumentCollection
	Out segment.Segment
}


func NewParser(groups int, Input *preprocessing.TkDocumentCollection){
	criterion := consts.Criterion
	length := len(criterion)
	step :=  length / groups

	parser := Parser{
		runeGroups: make([]rune, 0),
		Input: Input,
		Out: segment.Segment{
			Fragments: make([]*segment.Fragment, 0),
		},
	}	
	for i := 1; i <= groups + 1; i++{
		parser.runeGroups = append(parser.runeGroups, consts.Criterion[min(groups * step, length) - 1])
	}
}

func(p *Parser) Serve() {
	for _, d := range p.Input.DocList{
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
			p.Out.Fragments[idx].AddPair(t, d)
		}
	}
}