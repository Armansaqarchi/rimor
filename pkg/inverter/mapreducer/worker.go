package mapreducer

import (
	"rimor/pkg/inverter/mapreducer/segment"
	preprocessing "rimor/pkg/preprocessing"
)

type IParser interface {
	Serve(Input *preprocessing.TkDocumentCollection)
}


type IInverter interface {
	Serve(Input *segment.Segment)
}
