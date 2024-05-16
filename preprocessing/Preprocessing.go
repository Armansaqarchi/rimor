package preprocessing

type preprocessingDTO struct {
	docObject document
	appliedSteps []preprocessingStep
}

type preprocessingStep interface {
	Process(dto preprocessingDTO) preprocessingDTO
}

type preprocessor struct {
	inputData documentCollection
	preprocessingSteps []preprocessingStep
}

func (instance *preprocessor) SetInputData(inputData documentCollection) {
	instance.inputData = inputData
}

func NewPreprocessor(inputData documentCollection,  steps ...[]preprocessingStep) preprocessor {
	preprocessor := preprocessor{}
	preprocessor.inputData = inputData
	if len(steps) == 0 {
		preprocessor.preprocessingSteps = make([]preprocessingStep, 0)
		return preprocessor
	} 
	preprocessor.preprocessingSteps = steps[0]
	
	return preprocessor
}

func (preprocessor *preprocessor) Process(dto preprocessingDTO) preprocessingDTO {
	for _, prepStep := range preprocessor.preprocessingSteps {
		dto = prepStep.Process(dto)
	}
	return dto
}
