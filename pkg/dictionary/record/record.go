package record


type Recorder interface {
	GetTerm() string
	GetDF() int
	GetTF(int) int
	GetWeight(int) float64
	GetLast() IPostingListElem
	AddToPosting(IPostingListElem)
	GetPostingList() IPostingListElem
}



type Record struct {
	Term string
	DF int // no need to keep weights in float64 which is expensive
	PostingList IPostingListElem // points to the first element in posting list
	Last IPostingListElem
}


func NewRecord(term string, DocID int64) Recorder{
	DocRef := NewPostingListElem(DocID, nil)
	return &Record{
		Term: term,
		DF: 0,
		PostingList: DocRef,
		Last: DocRef,
	}
}


func (r *Record) GetTerm() string{
	return r.Term
}


func (r *Record) GetDF() int{
	return r.DF
}


func (r *Record) GetWeight(int) float64 {
	return 0.0
}

func (r *Record) GetTF(int) int {
	return 0
}

func (r *Record) GetLast() IPostingListElem {
	return r.Last
}


func (r *Record) AddToPosting(elm IPostingListElem) {
	if r.PostingList == nil {
		r.PostingList = elm
		r.Last = elm
	}
	r.Last.SetNextElem(elm)
	r.Last = elm
}

func (r *Record) GetPostingList() IPostingListElem{
	return r.PostingList
}