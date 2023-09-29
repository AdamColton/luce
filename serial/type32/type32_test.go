package type32_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/stretchr/testify/assert"
)

func TestErrTypeNotFound(t *testing.T) {
	intTp := reflector.Type[int]()
	assert.Equal(t, "type32: type int was not found", type32.ErrTypeNotFound{intTp}.Error())

}

type person struct {
	Name string
	Age  int
}

func (*person) TypeID32() uint32 {
	return 12345
}

type strSlice []string

func (strSlice) TypeID32() uint32 {
	return 67890
}
