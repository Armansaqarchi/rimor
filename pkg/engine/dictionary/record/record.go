package record

type Recorder interface {
	GetTerm() string
	GetDF() int64
	GetWeight(int) float64
	GetLast() IPostingElem
	AddToPosting(IPostingElem)
	IncreaseDF()
	GetPostingList() IPostingElem
	CreateChampions()
	GetChampion() IPostingElem
}

type Record struct {
	Term        string
	DF          int64        // no need to keep weights in float64 which is expensive
	PostingList IPostingElem // points to the first element in posting list
	Last        IPostingElem
	ChampionList IPostingElem
	ChamptionNum int
}

func NewRecord(term string, DocID int64) Recorder {
	DocRef := NewPostingListElem(DocID, nil)
	return &Record{
		Term:        term,
		DF:          0,
		PostingList: DocRef,
		Last:        DocRef,
	}
}

func (r *Record) GetTerm() string {
	return r.Term
}

func (r *Record) GetDF() int64 {
	return r.DF
}

func (r *Record) GetWeight(int) float64 {
	return 0.0
}

func (r *Record) GetLast() IPostingElem {
	return r.Last
}

func (r *Record) AddToPosting(elm IPostingElem) {
	if r.PostingList == nil {
		r.PostingList = elm
		r.Last = elm
	}
	r.Last.SetNextElem(elm)
	r.Last = elm
}

func (r *Record) GetPostingList() IPostingElem {
	return r.PostingList
}

func (r *Record) IncreaseDF() {
	r.DF++
}

func (r *Record) CreateChampions(){
	championList := []IPostingElem{}
	insert := func(e IPostingElem) {
		if e.GetTF() < championList[0].GetTF(){
			return
		}
		championList[0] = e
		idx := 0
		for idx < r.ChamptionNum-1 && (championList[idx].GetTF() < championList[idx+1].GetTF()) {
			tmp := championList[idx+1]
			championList[idx] = championList[idx+1]
			championList[idx+1] = tmp
		}
	}
	curr := r.PostingList
	for curr != nil {
		insert(curr)
		curr = curr.GetNextElem()
	}

	for i := r.ChamptionNum; i > 0 ; i++ {
		championList[i].SetNextElem(championList[i-1])
	}

	r.ChampionList = championList[r.ChamptionNum-1]
}

func (r *Record) GetChampion() IPostingElem {
	return r.ChampionList
}