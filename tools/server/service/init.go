package service

import (
	"encoding/gob"

	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/util/lfile"
)

var tm = type32.NewTypeMap()
var OS lfile.FSFileReader = lfile.OSRepository{}

// Register types with both gob and the typemap.
func Register(zeroValues ...type32.TypeIDer32) {
	for _, z := range zeroValues {
		gob.Register(z)
	}
	tm.RegisterType32s(zeroValues...)
}

func init() {
	Register(
		(*Request)(nil),
		(*Response)(nil),
		SocketOpened{},
	)
}
