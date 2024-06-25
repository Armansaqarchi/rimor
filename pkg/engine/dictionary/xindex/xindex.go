package xindex

import (
	"fmt"
	"rimor/pkg/engine/dictionary/record"
	"rimor/pkg/utils/errors"
)


type XIndexer interface{
	GetRecord(string) record.Recorder
	GetRecords() []record.Recorder
	GetDocumentsNum() int
	BinarySearchRecord() (record.Recorder, error)
}

type Xindex struct {
	DocLengths []int64
	Records []record.Recorder
	DocNum int64
	ChampionNum int
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

func (x *Xindex) BinarySearchRecord(t string) (record.Recorder, error){
	if x.Records == nil {
		return nil, fmt.Errorf("failed to get the target record, Records list are empty")
	}
	s, e := 0, len(x.Records)-1


	for s <= e {
		if s == e{
			if t == x.Records[s].GetTerm(){
				return x.Records[s], nil
			} else {
				return nil, errors.NewRecordNotFound(nil)
			}
		} else {
			mid := (s + e) / 2
			if t >= x.Records[mid].GetTerm(){
				s = mid+1
			} else {
				e = mid-1
			}
		}
	}

	return nil, fmt.Errorf("something went wrong while doing binary search")
} 









