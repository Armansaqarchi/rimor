package mapreducer

import (
	consts "rimor/pkg/consts"
	segment "rimor/pkg/engine/inverter/mapreducer/segment"
	preprocessing "rimor/pkg/engine/preprocessing"
	"unicode/utf8"
)

type Parser struct {
	runeGroups []rune
	Out        *segment.Segment
	
}

func NewParser(groups int, Out *segment.Segment) Parser {
	groups--
	criterion := consts.Criterion
	length := len(criterion)
	step := length / groups

	parser := Parser{
		runeGroups: make([]rune, 0),
		Out:        Out,
	}
	for i := 1; i <= groups+1; i++ {
		parser.runeGroups = append(parser.runeGroups, consts.Criterion[min(i*step, length)-1])
	}

	return parser
}

func (p *Parser) Serve(Input *preprocessing.TkDocumentCollection) {
	for _, doc := range Input.DocList {
		for x, token := range doc.TokenzedDocContent {
			p.AddTokenToFragment(token, int64(doc.Id), x)
		}
	}

}

func (p *Parser) AddTokenToFragment(token string, docId int64, x int) {
	if len(token) <= 0 {
		return
	}
	for idx, r := range p.runeGroups {
		t, _ := utf8.DecodeRuneInString(token)
		if t <= r {
			p.Out.Fragments[idx].AddPair(&segment.Pair{
				Term: token,
				Doc:  docId,
				Num: int64(x),
			})
			return
		}
	}
}
