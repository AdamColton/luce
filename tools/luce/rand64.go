package luce

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/adamcolton/luce/util/handler"
)

type RandBase64Command struct {
	Disabled bool
}

type Rand64Req struct {
	N int
}

type Rand64Resp struct {
	R string
}

func (rbh *RandBase64Command) RandBase64Handler(r *Rand64Req) *Rand64Resp {
	b := make([]byte, r.N)
	rand.Read(b)
	return &Rand64Resp{
		R: base64.URLEncoding.EncodeToString(b),
	}
}

func (rbh *RandBase64Command) RandBase64Usage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage:    "Creates a base64 encoded string N bytes long",
		Alias:    "r64",
		Disabled: rbh.Disabled,
	}
}
