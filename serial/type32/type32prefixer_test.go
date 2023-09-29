package type32_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/rye"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/stretchr/testify/assert"
)

func TestType32Prefixer(t *testing.T) {
	p := &person{
		Name: "Adam",
		Age:  39,
	}
	t32p := type32.Type32Prefixer{}
	got, err := t32p.PrefixInterfaceType(p, nil)
	assert.NoError(t, err)

	expected := make([]byte, 4)
	rye.Serialize.Uint32(expected, p.TypeID32())

	assert.Equal(t, expected, got)

	_, err = t32p.PrefixInterfaceType(123, nil)
	assert.Equal(t, type32.ErrTypeNotFound, err)

	s := t32p.Serializer(nil)
	assert.Equal(t, s.InterfaceTypePrefixer, t32p)
}
