package mapreducer

import (
	"rimor/pkg/engine/inverter/mapreducer/segment"
	preprocessing "rimor/pkg/engine/preprocessing"
)

type IParser interface {
	Serve(Input *preprocessing.TkDocumentCollection)
}


type IInverter interface {
	Serve(Input *segment.Segment)
}
