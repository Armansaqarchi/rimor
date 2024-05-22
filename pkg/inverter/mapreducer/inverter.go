package mapreducer

 import (
	"rimor/pkg/inverter/mapreducer/segment"
	"rimor/pkg/dictionary/xindex"
 )

type Inverter struct {
	Number int
	Input *segment.Segment
	Out xindex.Xindex
}


func (inv *Inverter) Serve() {

}