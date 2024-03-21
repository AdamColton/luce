package luce

import (
	"crypto/rand"
	"encoding/base64"
)

type Rand64Req struct {
	N int
}

type Rand64Resp struct {
	R string
}

func (l *Luce) RandBase64Handler(r *Rand64Req) *Rand64Resp {
	b := make([]byte, r.N)
	rand.Read(b)
	return &Rand64Resp{
		R: base64.URLEncoding.EncodeToString(b),
	}
}

func (l *Luce) RandBase64Usage() string {
	return "Creates a base64 encoded string N bytes long"
}
