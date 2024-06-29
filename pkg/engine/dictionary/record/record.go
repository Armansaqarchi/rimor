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
	Term         string
	DF           int64        // no need to keep weights in float64 which is expensive
	PostingList  IPostingElem // points to the first element in posting list
	Last         IPostingElem
	ChampionList IPostingElem
	ChampionNum  int
}

func NewRecord(term string, DocID int64, Num int64) Recorder {
	DocRef := NewPostingListElem(DocID, nil)
	defaultChampionNum := 1024
	r := &Record{
		Term:        term,
		DF:          0,
		PostingList: DocRef,
		Last:        DocRef,
		ChampionNum: defaultChampionNum,
	}
	r.GetLast().AddPosition(Num)
	return r
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

func (r *Record) CreateChampions() {
	var championList []IPostingElem = nil
	insert := func(e IPostingElem) bool {
		if championList == nil {
			championList = make([]IPostingElem, r.ChampionNum)
			championList[0] = e
			return true
		} else {
			for i := 0; i < r.ChampionNum; i++ {
				if championList[i] == nil {
					championList[i] = e
					return true
				} else if e.GetTF() > championList[i].GetTF() {
					for j := r.ChampionNum - 1; j > i; j-- {
						championList[j] = championList[j-1]
					}
					championList[i] = e
					return true
				}
			}
			return false
		}
	}
	curr := r.PostingList
	insertCount := 0
	for curr != nil {
		toBeInserted := PostingElem{
			DocID:    curr.GetDocID(),
			TF:       curr.GetTF(),
			NextElem: nil,
		}
		if insert(&toBeInserted) {
			insertCount++
		}
		curr = curr.GetNextElem()
	}

	insertCount = min(insertCount, r.ChampionNum)
	for i := insertCount - 1; i > 0; i-- {
		championList[i].SetNextElem(championList[i-1])
	}

	r.ChampionList = championList[insertCount-1]
}

func (r *Record) GetChampion() IPostingElem {
	return r.ChampionList
}
