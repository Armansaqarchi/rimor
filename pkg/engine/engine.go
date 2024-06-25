package engine

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"log"
	"os"
	consts "rimor/pkg/consts"
	"rimor/pkg/engine/dictionary/xindex"
	"rimor/pkg/engine/inverter/mapreducer"
	MReduce "rimor/pkg/engine/inverter/mapreducer"
	preprocessing "rimor/pkg/engine/preprocessing"
	tokenizer "rimor/pkg/engine/preprocessing/tokenizer"
	"rimor/pkg/scoring"
)


type Engine struct {
	DocumentCollection *preprocessing.DocumentCollection
	Preprocessor preprocessing.Preprocessor
	Tokenizer *tokenizer.WordTokenizer
	Constructor *mapreducer.Master
	Index		*xindex.Xindex
	K 			int
}


func readDocumentCollection(documentCollectionPath string) (preprocessing.DocumentCollection, error) {
	documentCollectionAsJsonBytes, err := os.ReadFile(documentCollectionPath)
	if err != nil {
		return preprocessing.DocumentCollection{}, err
	}


	var documentCollection preprocessing.DocumentCollection
	if err := json.Unmarshal(documentCollectionAsJsonBytes, &documentCollection); err != nil {
		return preprocessing.DocumentCollection{}, err
	}

	return documentCollection, nil
}

func NewEngine() *Engine{
	docCollection, err := readDocumentCollection(consts.COLLECTION_PATH)
	if err != nil {
		fmt.Println(err)
	}

	TokenizedCollection := preprocessing.TkDocumentCollection{
		DocList: make([]preprocessing.TkDocument, 0),
	}
	tokenizer, err := tokenizer.NewWordTokenizer(tokenizer.WORDS_PATH, tokenizer.VERBS_PATH, false, false, false, false, false, false, false, false)
	if err != nil {
		log.Fatalf("failed to instantiate the tokenizer, err : %s", err.Error())
	}
	for _, col := range docCollection.DocList {
		TokenizedCollection.DocList = append(TokenizedCollection.DocList, preprocessing.TkDocument{
			Id: col.ID,
			TokenzedDocContent : tokenizer.Tokenize(col.Content),
			DocUrl: col.Url,
		}) 
	}

	MapReducer := MReduce.NewMaster(8, len(TokenizedCollection.DocList)/4, 30)
	indx := MapReducer.MapReduce(TokenizedCollection)

	engine := Engine{
		DocumentCollection: &docCollection,
		Preprocessor: preprocessing.Preprocessor{},
		Tokenizer: tokenizer,
		Constructor: MapReducer,
		Index: indx,
		K: 30,
	}

	return &engine
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


	tokenizedQuery := e.Tokenizer.Tokenize(tq)
	queryLen := len(tokenizedQuery)
	queryTermMap := make(map[string] int8)


	vectorizedQuery := make([]VectorElem, queryLen)

	for _, token := range tokenizedQuery {
		val, contains := queryTermMap[token]
		if contains {
			queryTermMap[token] = val +1
		} else {
			queryTermMap[token] = 0 
		}	
	}

	for idx, token := range tokenizedQuery {
		tokenTF := queryTermMap[token]
		vectorizedQuery[idx] = VectorElem{
			Term: token,
			Value: int64(tokenTF),
		}
	}

	// preprocessing steps for query
	q := Query{
		Vector: vectorizedQuery,
	} // this has to be populated after preprocessing step on text query

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