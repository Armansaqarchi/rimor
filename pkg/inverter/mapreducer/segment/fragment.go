package segment


type Fragment struct {
	Pairs []*Pair
}

func (f *Fragment) AddPair(pair *Pair){
	f.Pairs = append(f.Pairs, pair)
}

// impl of internal go sort pkg
func (f Fragment) Len() int {
	return len(f.Pairs)
}
func (f Fragment) Less(i, j int) bool {
	return (*f.Pairs[i]).Term < (*f.Pairs[j]).Term
}
func (f Fragment) Swap(i, j int) {
	f.Pairs[i], f.Pairs[j] = f.Pairs[j], f.Pairs[i]
}
// end of go sort pkg impl


type Pair struct {
	Term string
	Doc  int64
}