package xindex


import (
	"rimor/pkg/dictionary/record"
)


type XIndexer interface{
	GetRecord(string) record.Recorder
	GetRecords() []record.Recorder
	GetDocumentsNum() int
	Sort()
}

type Xindex struct {
	Records []record.Recorder
	DocNum int
	Sorted bool
}


func (x *Xindex) GetRecord(string) record.Recorder {
	return nil
}

func (x *Xindex) GetRecords(string) []record.Recorder {
	return []record.Recorder{}
}

func (x *Xindex) GetDocumentsNum() int {
	return 0
}

func (x *Xindex) Sort() int {
	return 0
}






