package scoring

type ScoreFN func(w_tq int, w_td int, dl float64) float64 

/* CosineScore
	w_tq represents the weight for term t in a query q
	w_td represents the weight for term t in a doc d
	dl represents the document length
*/
func CosineScore(w_tq, w_td float64, dl int64) float64{
	return w_td * w_tq / float64(dl)
}