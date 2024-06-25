package record

type IPostingElem interface {
	GetTF() int64
	GetDocID() int64
	GetNextElem() IPostingElem
	SetNextElem(IPostingElem)
	IncreaseTF()
}

type PostingElem struct {
	DocID    int64
	TF       int64
	NextElem IPostingElem
}

func NewPostingListElem(docID int64, next IPostingElem) IPostingElem {
	return &PostingElem{
		TF:       1,
		DocID:    docID,
		NextElem: next,
	}
}

func (ple *PostingElem) GetTF() int64 {
	return ple.TF
}

func (ple *PostingElem) GetDocID() int64 {
	return ple.DocID
}

func (ple *PostingElem) GetNextElem() IPostingElem {
	return ple.NextElem
}

func (ple *PostingElem) SetNextElem(elm IPostingElem) {
	ple.NextElem = elm
}

func (ple *PostingElem) IncreaseTF() {
	ple.TF++
}
