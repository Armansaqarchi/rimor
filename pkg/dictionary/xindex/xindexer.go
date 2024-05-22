package dictionary


import (
	"rimor/pkg/dictionary/record"
)


type XIndexer interface{
	GetRecord(string) record.Recorder
}



