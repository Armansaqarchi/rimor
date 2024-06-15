package scoring

import "math"



func TF_IDF(TF, IDF, N int64) float64{
	if TF == 0 {
		return 0.0
	} else {
		return (1 + math.Log(float64(TF))) * math.Log(float64(N)/float64(IDF))
	}
}

