package errors

import "fmt"


type BaseErr struct{
	m any
}


func (r BaseErr) Error() string {
	return fmt.Sprintf("%v", r.m)
}

type RecordNotFound struct {
	BaseErr
}


func NewRecordNotFound(m any) RecordNotFound{
	return RecordNotFound{
		BaseErr: BaseErr{
			m: m,
		},
	}
}





