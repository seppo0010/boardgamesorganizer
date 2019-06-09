package time

import "time"

type Real struct {
}

func NewReal() *Real {
	return &Real{}
}

func (r *Real) Now() time.Time {
	return time.Now()
}
