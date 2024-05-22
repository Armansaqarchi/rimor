package record


type IPostingListElem interface {
	GetTF() int
	GetDocID() int
	GetNextElemt() *IPostingListElem
}


type PostingListElem struct {
	DocID 	int
	TF  	int
	NextElem IPostingListElem
}


func (ple *PostingListElem) GetDF() int{
	return ple.TF
}

func (ple *PostingListElem) GetDocID() int {
	return ple.DocID
}

func (ple *PostingListElem) GetNextElemt() IPostingListElem{
	return ple.NextElem
}




