package record


type IPostingListElem interface {
	GetTF() int64
	GetDocID() int64
	GetNextElem() IPostingListElem
	SetNextElem(IPostingListElem)
	IncreaseTF()
}


type PostingListElem struct {
	DocID 	int64
	TF  	int64
	NextElem IPostingListElem
}


func NewPostingListElem(docID int64, next IPostingListElem) IPostingListElem {
	return &PostingListElem{
		TF: 1,
		DocID: docID,
		NextElem: next,
	}
}


func (ple *PostingListElem) GetTF() int64{
	return ple.TF
}

func (ple *PostingListElem) GetDocID() int64 {
	return ple.DocID
}

func (ple *PostingListElem) GetNextElem() IPostingListElem{
	return ple.NextElem
}


func (ple *PostingListElem) SetNextElem(elm IPostingListElem) {
	ple.NextElem = elm
}

func (ple *PostingListElem) IncreaseTF() {
	ple.TF++
}
