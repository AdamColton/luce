package luce

import (
	"math/rand"

	"github.com/adamcolton/luce/util/handler"
)

type RandCommand struct {
	Disabled bool
}

type RandReq struct {
	N    int64
	Bits bool `prompt:"Y/n"`
}

type RandResp struct {
	R int64
}

func (rh *RandCommand) RandHandler(r *RandReq) *RandResp {
	if r.Bits {
		r.N = 1 << r.N
	}
	return &RandResp{
		R: rand.Int63n(r.N),
	}
}

func (rh *RandCommand) RandUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage:    "Create a random int. N determines the max, if Bits is true, N is the number of bits.",
		Disabled: rh.Disabled,
	}
}
