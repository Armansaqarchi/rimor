package preprocessing

type DocumentCollection struct {
	DocList []Document
}

type Document struct {
	Id         int64  `json:"id"`
	DocUrl     string `json:"url"`
	DocContent string `json:"text"`
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
	Id         int64 
	DocUrl     string
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