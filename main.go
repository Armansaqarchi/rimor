package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	MReduce "rimor/pkg/inverter/mapreducer"
	"rimor/pkg/preprocessing"
	tokenizer "rimor/pkg/preprocessing/tokenizer"
)

const COLLECTION_PATH = "./data/news.json"



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

func main() {
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
	index := MapReducer.MapReduce(TokenizedCollection)
	for _, p := range index.Records {
		if p.GetTerm()[0] < 127 {
			continue
		}
		fmt.Print(p.GetTerm())
		current := p.GetPostingList()
		for current != nil {
			fmt.Printf("%d ", current.GetDocID())
			current = current.GetNextElem()
		}
		fmt.Print("\n")
	}

}
