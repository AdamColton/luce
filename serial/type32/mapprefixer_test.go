package type32_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/rye"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/stretchr/testify/assert"
)

func TestMapPrefixer(t *testing.T) {
	strTp := reflector.Type[string]()
	var mp type32.MapPrefixer
	_, err := mp.PrefixReflectType(strTp, nil)
	assert.Equal(t, type32.ErrTypeNotFound{strTp}, err)

	mp = make(type32.MapPrefixer)
	mp[strTp] = 1234

	expected := make([]byte, 4)
	rye.Serialize.Uint32(expected, 1234)

	b, err := mp.PrefixReflectType(strTp, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, b)

	intTp := reflector.Type[int]()
	check, err := mp.PrefixReflectType(intTp, b)
	assert.Equal(t, b, check)
	assert.Equal(t, type32.ErrTypeNotFound{intTp}, err)

	s := mp.Serializer(nil)
	b, err = s.PrefixInterfaceType("test", b[:0])
	assert.NoError(t, err)
	assert.Equal(t, expected, b)
}

func TestErrTypeNotFound(t *testing.T) {
	intTp := reflector.Type[int]()
	assert.Equal(t, "Type int was not found", type32.ErrTypeNotFound{intTp}.Error())

}
