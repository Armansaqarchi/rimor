package preprocessing

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type DocumentCollection struct {
	DocList []Document
}

type Document struct {
	ID       int64 
	Title    string 
	Content  string 
	Date     string     // change the type to time.time if necessary
	Url      string 
	Tags	[]string 
	Category string 
}

type Documents struct {
	DocumentList []Document
}

func (ds *Documents) UnmarshalJSON(d []byte) error {

	var m map[string]Document
	if err := json.Unmarshal(d, &m); err != nil {
		log.Printf("failed to parse raw string into map, err : %s", err.Error())
	}
	for k, v := range m {
		id, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to convert string id into integer id, err : %s", err.Error())
		}
		ds.DocumentList = append(ds.DocumentList, Document{
			ID: id,
			Title: v.Title,
			Content: v.Content,
			Date: v.Date,
			Url: v.Url,
			Tags: v.Tags,
			Category: v.Category,

		})
	}

	return nil
}

type PreprocessingStep interface {
	Process(dto Document) Document
}

type Preprocessor struct {
	inputData          DocumentCollection
	preprocessingSteps []PreprocessingStep
}

func (instance *Preprocessor) SetInputData(inputData DocumentCollection) {
	instance.inputData = inputData
}

func NewPreprocessor(inputData DocumentCollection, steps ...[]PreprocessingStep) Preprocessor {
	preprocessor := Preprocessor{}
	preprocessor.inputData = inputData
	if len(steps) == 0 {
		preprocessor.preprocessingSteps = make([]PreprocessingStep, 0)
		return preprocessor
	}
	preprocessor.preprocessingSteps = steps[0]

	return preprocessor
}

func (preprocessor *Preprocessor) Process() <-chan Document {
	preProcessResChannel := make(chan Document, 4)
	for _, document := range preprocessor.inputData.DocList {
		defer close(preProcessResChannel)
		preProcessResChannel <- preprocessor.applyAllPreprocessingSteps(document)
	}
	return preProcessResChannel
}

func (preprocessor *Preprocessor) applyAllPreprocessingSteps(document Document) Document {
	var processed Document = document
	for _, step := range preprocessor.preprocessingSteps {
		processed = step.Process(processed)
	}
	return processed
}

type TkDocument struct {
	Id                 int64
	DocUrl             string
	TokenzedDocContent []string
}

type TkDocumentCollection struct {
	DocList []TkDocument
}

type Tokenizer struct {
}

func (tokenizer *Tokenizer) Tokenize(documentCollection DocumentCollection) TkDocumentCollection {
	return TkDocumentCollection{}
}
