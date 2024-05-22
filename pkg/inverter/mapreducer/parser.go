package mapreducer

import (
	"rimor/pkg/inverter/mapreducer/segment"
)

type Split struct {

}

// the above struct represents the input data, containing information for each documents and its related tokens

type Parser struct {
	ParseGroups int
	Split Split
	SegmentOut segment.Segment
}


func(p *Parser) Serve() {
	
}