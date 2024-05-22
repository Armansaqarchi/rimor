package record


type Recorder interface {
	GetTerm() string
	GetIDF() int
	GetTermFreq() int
	GetDfList() []int
	GetDfNonZeroList() []int
}



type Record struct {
	Term string
	TermFreq int
	IDF int // no need to keep weights in float64 which is expensive
	PostingList *PostingListElem // points to the first element in posting list
}


func (r *Record) GetTerm() string{
	return r.Term
}


func (r *Record) GetIDF() int{
	return r.IDF
}

func (r *Record) GetDFList() []int {
	return []int{}
}

func (r *Record) GetDfNonZeroList() []int {
	return []int{}
}