package record


type IPostingListElem interface {
	GetDF() int
	GetDocID() int
	GetNextElemt() *IPostingListElem
}


type PostingListElem struct {
	DocID 	int
	DF  	int
	NextElem IPostingListElem
}


func (ple *PostingListElem) GetDF() int{
	return ple.DF
}

func (ple *PostingListElem) GetDocID() int {
	return ple.DocID
}

func (ple *PostingListElem) GetNextElemt() IPostingListElem{
	return ple.NextElem
}




