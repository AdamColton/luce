package type32_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/serial/type32"
	"github.com/stretchr/testify/assert"
)

func TestPrefixIsLittleEndian(t *testing.T) {
	typ := reflect.TypeOf("")
	p := type32.MapPrefixer{typ: 0x01020304}

	got, err := p.PrefixReflectType(typ, []byte{9})
	assert.NoError(t, err)
	assert.Equal(t, []byte{9, 4, 3, 2, 1}, got)
}
