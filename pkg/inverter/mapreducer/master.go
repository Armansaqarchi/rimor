package mapreducer

import "github.com/Jeffail/tunny"



type Master struct {
	Splits int
	Pool tunny.Pool
}


// func NewMaster(splits int, routines int) *Master{
// 	return &Master{
// 		Splits: splits,
// 		Pool: *tunny.NewFunc(),
// 	}
// } 