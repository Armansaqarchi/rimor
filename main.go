package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"rimor/pkg/engine/dictionary/xindex"
	index "rimor/pkg/engine/dictionary/xindex"
	MReduce "rimor/pkg/engine/inverter/mapreducer"
	"rimor/pkg/engine/preprocessing"
	tokenizer "rimor/pkg/engine/preprocessing/tokenizer"
)

const COLLECTION_PATH = "./data/news.json"
const INDEX_FILE_PATH = "./data/tmp.json"


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


var indx *index.Xindex

func main() {
	data_present := flag.Bool("present", false, "this is mainly used for restoring of index when index has previously been calculated")
	store := flag.Bool("storedata", false, "when this flag is set, index data will be stored in file")
	flag.Parse()

	if *data_present{
		indx = ImportIndexToFile()
	} else {
		docCollection, err := readDocumentCollection(COLLECTION_PATH)
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
		indx = MapReducer.MapReduce(TokenizedCollection)
	}

	if *store && !(*data_present) {
		ExportIndexToFile(indx)
	}

	for _, p := range indx.Records {
		fmt.Print(p.GetTerm())
		current := p.GetPostingList()
		for current != nil {
			fmt.Printf("%d ", current.GetDocID())
			current = current.GetNextElem()
		}
		fmt.Print("\n")
	}

}


func ExportIndexToFile(index *index.Xindex) error{
	f, err := os.Create(INDEX_FILE_PATH)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	enc.Encode(index)

	return nil
}

func ImportIndexToFile() *xindex.Xindex{
	f, err := os.Open(INDEX_FILE_PATH)
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
