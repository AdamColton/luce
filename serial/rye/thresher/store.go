package thresher

import (
	"github.com/adamcolton/luce/ds/lmap"
)

var store = lmap.Map[string, []byte]{}
