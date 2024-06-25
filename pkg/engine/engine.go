package engine

import (
	"container/heap"
	"fmt"
	"rimor/pkg/engine/dictionary/xindex"
	"rimor/pkg/engine/inverter/mapreducer"
	preprocessing "rimor/pkg/engine/preprocessing"
	"rimor/pkg/scoring"

)


type Engine struct {
	DocumentCollection preprocessing.DocumentCollection
	Preprocessor preprocessing.Preprocessor
	Constructor mapreducer.Master
	Index		xindex.Xindex
	K 			int
}


type VectorElem struct {
	Term string
	Value int64
}

type Query struct {
	Vector []VectorElem
}


func (e *Engine) Score(q Query) ([]float64, error){
	scores := make([]float64, e.Index.DocNum)
	for _, t := range q.Vector{
		r, err := e.Index.BinarySearchRecord(t.Term)
		if err != nil {
			return nil, err
		}
		p := r.GetPostingList()
		for p != nil {
			Wtd := scoring.TF_IDF(p.GetTF(), r.GetDF(), e.Index.DocNum)
			Wtq := scoring.TF_IDF(t.Value, r.GetDF(), e.Index.DocNum)
			scores[p.GetDocID()] += scoring.CosineScore(Wtq, Wtd, e.Index.DocLengths[p.GetDocID()])
			p = p.GetNextElem()
		}
	}

	return scores, nil
}


func (e *Engine) Query(tq string)(*preprocessing.DocumentCollection, error) {

	// preprocessing steps for query
	q := Query{} // this has to be populated after preprocessing step on text query

	scores, err := e.Score(q)
	if err != nil {
		return nil, fmt.Errorf("failed to process query, err : %s", err.Error())
	}
	sh := DocumentHeap{}
	for id, score := range scores {
		sh = append(sh, DocumentScore{
			DocID: id,
			Score: score,
		})
	}
	heap.Init(&sh)

	DocCollection := preprocessing.DocumentCollection{}

	for i := 0; i < e.K; i++ {
		ds , ok:= heap.Pop(&sh).(DocumentScore)
		if !ok {
			return nil, fmt.Errorf("something went wrong while evaluation of documents")
		}
		DocCollection.DocList = append(DocCollection.DocList, e.DocumentCollection.DocList[ds.DocID])
	}
	return &DocCollection, nil

}


type DocumentScore struct {
	DocID 	int
	Score   float64
}

type DocumentHeap []DocumentScore


func (d DocumentHeap) Len() int{
	return len(d)
}

func (d DocumentHeap) Less(first , second int) bool{
	return d[first].Score > d[second].Score // this is reversed due to achievement of max heap
}

func (d DocumentHeap) Swap(first, second int) {
	d[first], d[second] = d[second], d[first]
}


func (d *DocumentHeap) Push(x any){
	*d = append(*d, x.(DocumentScore))
}

func (d *DocumentHeap) Pop() any{
	old := *d
	n := len(old)
	p := old[n-1]
	*d = old[0: n-1]
	return p
}