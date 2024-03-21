package luce

import "math/rand"

type RandReq struct {
	N    int64
	Bits bool `prompt:"Y/n"`
}

type RandResp struct {
	R int64
}

func (l *Luce) RandHandler(r *RandReq) *RandResp {
	if r.Bits {
		r.N = 1 << r.N
	}
	return &RandResp{
		R: rand.Int63n(r.N),
	}
}

func (l *Luce) RandUsage() string {
	return "Create a random int. N determines the max, if Bits is true, N is the number of bits."
}
