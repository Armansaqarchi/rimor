package main

import (
	"encoding/json"
	"fmt"
	"os"
	"rimor/preprocessing"
)

const COLLECTION_PATH = "collection/test_collection.json"

func readDocumentCollection(documentCollectionPath string) (preprocessing.DocumentCollection, error) {
	documentCollectionAsJsonBytes, err := os.ReadFile(documentCollectionPath)
	if err != nil {
		return preprocessing.DocumentCollection{}, err
	}


	var documentCollection preprocessing.DocumentCollection
	if err := json.Unmarshal(documentCollectionAsJsonBytes, &documentCollection.DocList); err != nil {
		return preprocessing.DocumentCollection{}, err
	}

	return documentCollection, nil
}

func main() {
	docCollection, err := readDocumentCollection(COLLECTION_PATH)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(docCollection)
}
