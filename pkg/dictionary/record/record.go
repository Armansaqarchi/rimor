package record


type Recorder interface {
	GetTerm() string
	GetDF() int
	GetTF(int) int
	GetWeight(int) float64
	GetTermFreq() int
	GetTFList() []int
	GetDfNonZeroList() []int
}



type Record struct {
	Term string
	TermFreq int
	DF int // no need to keep weights in float64 which is expensive
	PostingList *PostingListElem // points to the first element in posting list
}


func (r *Record) GetTerm() string{
	return r.Term
}


func (r *Record) GetDF() int{
	return r.DF
}


func (r *Record) GetWeight() float64 {
	return 0.0
}

func (r *Record) GetTF(int) int {
	return 0
}

func (r *Record) GetDFList() []int {
	return []int{}
}

func (r *Record) GetDfNonZeroList() []int {
	return []int{}
}