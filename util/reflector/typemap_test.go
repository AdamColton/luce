package reflector_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

func TestTypeMap(t *testing.T) {
	tm := reflector.TypeMap{}
	reflector.TMAdd[string]("Name", tm)
	reflector.TMEmbed[slice.Slice[string]](tm)

	assert.Equal(t, ltype.String, tm["Name"])
	assert.Equal(t, reflector.Type[slice.Slice[string]](), tm["Slice"])

	tc := reflector.TypeCollection{
		ltype.String: tm,
	}
	assert.Equal(t, tm, reflector.TMGet[string](tc))
}
