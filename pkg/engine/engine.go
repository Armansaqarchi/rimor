package engine

import (
	"container/heap"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	consts "rimor/pkg/consts"
	"rimor/pkg/engine/dictionary/record"
	"rimor/pkg/engine/dictionary/xindex"
	"rimor/pkg/engine/inverter/mapreducer"
	MReduce "rimor/pkg/engine/inverter/mapreducer"
	preprocessing "rimor/pkg/engine/preprocessing"
	tokenizer "rimor/pkg/engine/preprocessing/tokenizer"
	"rimor/pkg/scoring"
	errors_util "rimor/pkg/utils/errors"
)

type Engine struct {
	DocumentCollection  *preprocessing.DocumentCollection
	Preprocessor        preprocessing.Preprocessor
	MostUsedWordRemover *preprocessing.MostUsedWordRemover
	Tokenizer           *tokenizer.WordTokenizer
	Constructor         *mapreducer.Master
	Index               *xindex.Xindex
	MaxResultCount      int
}

func readDocumentCollection(documentCollectionPath string) (preprocessing.DocumentCollection, error) {
	docFile, err := os.Open(documentCollectionPath)
	if err != nil {
		return preprocessing.DocumentCollection{}, err
	}
	var documentCollection preprocessing.DocumentCollection

	dec := json.NewDecoder(docFile)
	if err := dec.Decode(&documentCollection); err != nil {
		return preprocessing.DocumentCollection{}, err
	}
	return documentCollection, nil
}

func ExportIndexToFile(index *xindex.Xindex) error{
	f, err := os.Create(consts.INDEX_FILE_PATH)
	if err != nil {
		log.Fatalln(err.Error())
	}
	enc := json.NewEncoder(f)
	if err := enc.Encode(index.Records); err != nil {
		log.Fatalln(err.Error())
	}
	return nil
}

func ImportIndexToFile() *xindex.Xindex{
	f, err := os.Open(consts.INDEX_FILE_PATH)
	if err != nil {
		log.Fatalf("failed to restore index, err : %s", err.Error())
	}
	dec := json.NewDecoder(f)
	indx := xindex.Xindex{}
	if err := dec.Decode(&indx); err != nil {
		log.Fatalf("failed to decode index data into object, err : %s", err.Error())
	}
	return &indx	
}

func NewEngine() *Engine {


	data_present := flag.Bool("present", false, "this is mainly used for restoring of index when index has previously been calculated")
	store := flag.Bool("storedata", false, "when this flag is set, index data will be stored in file")
	flag.Parse()




	// logger.Info().Msg("reading documents from file...")
	docCollection, err := readDocumentCollection(consts.COLLECTION_PATH)
	if err != nil {
		fmt.Println(err)
	}

	// logger.Info().Msg("tokenizing documents...")

	TokenizedCollection := preprocessing.TkDocumentCollection{
		DocList: make([]preprocessing.TkDocument, 0),
	}

	
	wordTokenizer, err := tokenizer.NewWordTokenizer(tokenizer.WORDS_PATH, tokenizer.VERBS_PATH, true, false, false, false, false, true, false, false)
	if err != nil {
		log.Fatalf("failed to instantiate the wordTokenizer, err : %s", err.Error())
	}
	var arabicPhrase = preprocessing.NewSpecialArabicPhraseNormalizer()
	var persianDigit = preprocessing.NewPersianDigitNormalizer()
	var unicodeRep = preprocessing.NewUnicodeReplacementPersianNormalizer()
	var punctuationRemover = preprocessing.NewPunctuationRemover()
	// logger.Info().Msg("preprocessing documents...")
	preprocessor := preprocessing.NewPreprocessor([]preprocessing.PreprocessingStep{
		&arabicPhrase,
		&persianDigit,
		&unicodeRep,
		&punctuationRemover,
	})

	for idx, col := range docCollection.DocList {
		col.Content = preprocessor.Process(col.Content)
		if idx == 1 {
			fmt.Print(col.Content)
		}
		tokenized := wordTokenizer.Tokenize(col.Content)
		if idx == 1 {
			fmt.Print(col.Content)
		}
		TokenizedCollection.DocList = append(TokenizedCollection.DocList, preprocessing.TkDocument{
			Id:                 col.ID,
			TokenzedDocContent: tokenized,
			DocUrl:             col.Url,
		})
	}


	// logger.Info().Msg("removing unused and stop words...")
	mostUsedWordRemover := preprocessing.NewMostUsedWordRemover()
	TokenizedCollection = mostUsedWordRemover.ProcessDocCollection(TokenizedCollection)
	fmt.Print(TokenizedCollection.DocList[0].TokenzedDocContent)
	// logger.Info().Msg("stemming persian words...")
	stemmer := preprocessing.NewPersianStemmer()
	for _, doc := range TokenizedCollection.DocList {
		for idx, token := range doc.TokenzedDocContent {
			doc.TokenzedDocContent[idx] = stemmer.Stem(token)
		}
	}
	fmt.Print(TokenizedCollection.DocList[0].TokenzedDocContent)
	// logger.Info().Msg("creating inverted index using MapReduce concurrent algorithm...")
	MapReducer := MReduce.NewMaster(8, len(TokenizedCollection.DocList)/4, 30)
	indx := MapReducer.CreateIndex(TokenizedCollection)

	if *store && !(*data_present) {
		fmt.Print("exporting...\n\n\n")
		ExportIndexToFile(indx)
	}

	// logger.Info().Msg("X index created successfully.\n\n")
	engine := Engine{
		DocumentCollection:  &docCollection,
		Preprocessor:        preprocessor,
		MostUsedWordRemover: &mostUsedWordRemover,
		Tokenizer:           wordTokenizer,
		Constructor:         MapReducer,
		Index:               indx,
		MaxResultCount:      10,
	}

	return &engine
}

type VectorElem struct {
	Term  string
	Value int64
}

type Query struct {
	Vector []VectorElem
}

func (e *Engine) Score(q Query, useChampions bool) ([]float64, error) {
	scores := make([]float64, e.Index.DocNum)
	for _, t := range q.Vector {
		r, err := e.Index.BinarySearchRecord(t.Term)
		if errors.Is(err, errors_util.RecordNotFound{}) {
			fmt.Print("term not found\n")
			continue
		}
		if err != nil {
			fmt.Print(err.Error())
			return nil, err
		}

		var p record.IPostingElem = nil
		if useChampions {
			p = r.GetChampion()
		} else {
			p = r.GetPostingList()
		}

		for p != nil {
			Wtd := scoring.TF_IDF(p.GetTF(), r.GetDF(), e.Index.DocNum)
			Wtq := scoring.TF_IDF(t.Value, r.GetDF(), e.Index.DocNum)
			scores[p.GetDocID()] += scoring.CosineScore(Wtq, Wtd, e.Index.DocLengths[p.GetDocID()])
			p = p.GetNextElem()
		}
	}

	return scores, nil
}

func (e *Engine) Query(tq string, useChampions bool) (*preprocessing.DocumentCollection, error) {

	fmt.Print("processing query...\n")
	tq = e.Preprocessor.Process(tq)
	tokenizedQuery := e.MostUsedWordRemover.ProcessQuery(e.Tokenizer.Tokenize(tq))

	queryTermMap := make(map[string]int8)
	fmt.Print("vectorizing query...\n")

	vectorizedQuery := []VectorElem{}

	for _, token := range tokenizedQuery {
		val, contains := queryTermMap[token]
		if contains {
			queryTermMap[token] = val + 1
		} else {
			queryTermMap[token] = 1
		}
	}

	for k, v := range queryTermMap {
		vectorizedQuery = append(vectorizedQuery, VectorElem{
			Term:  k,
			Value: int64(v),
		})
	}

	q := Query{
		Vector: vectorizedQuery,
	}
	scores, err := e.Score(q, useChampions)
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

	for i := 0; i < e.MaxResultCount; i++ {
		ds, ok := heap.Pop(&sh).(DocumentScore)
		if !ok {
			return nil, fmt.Errorf("something went wrong while evaluating documents")
		}
		DocCollection.DocList = append(DocCollection.DocList, e.DocumentCollection.DocList[ds.DocID])
	}
	return &DocCollection, nil

}

type DocumentScore struct {
	DocID int
	Score float64
}

type DocumentHeap []DocumentScore

func (d DocumentHeap) Len() int {
	return len(d)
}

func (d DocumentHeap) Less(first, second int) bool {
	return d[first].Score > d[second].Score // this is reversed due to achievement of max heap
}

func (d DocumentHeap) Swap(first, second int) {
	d[first], d[second] = d[second], d[first]
}

func (d *DocumentHeap) Push(x any) {
	*d = append(*d, x.(DocumentScore))
}

func (d *DocumentHeap) Pop() any {
	old := *d
	n := len(old)
	p := old[n-1]
	*d = old[0 : n-1]
	return p
}
