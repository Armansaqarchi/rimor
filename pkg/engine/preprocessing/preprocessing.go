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

type PreprocessingStep interface {
	Process(dto string) string
}


func (ds *DocumentCollection) UnmarshalJSON(d []byte) error {

	var m map[string]Document
	if err := json.Unmarshal(d, &m); err != nil {
		log.Printf("failed to parse raw string into map, err : %s", err.Error())
	}
	for i := 0; i < len(m); i++ {
		id := strconv.Itoa(i)
		v, ok := m[id]
		if !ok {
			return fmt.Errorf("failed to get the docuemnts")
		}
		ds.DocList = append(ds.DocList, Document{
			ID: int64(i),
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

type Preprocessor struct {
	preprocessingSteps []PreprocessingStep
}

func NewPreprocessor(steps ...[]PreprocessingStep) Preprocessor {
	preprocessor := Preprocessor{}
	if len(steps) == 0 {
		preprocessor.preprocessingSteps = make([]PreprocessingStep, 0)
		return preprocessor
	}
	preprocessor.preprocessingSteps = steps[0]

	return preprocessor
}

func (preprocessor *Preprocessor) Process(text string) string {
	return preprocessor.applyAllPreprocessingSteps(text)
}

func (preprocessor *Preprocessor) applyAllPreprocessingSteps(text string) string {
	var processed string = text
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
