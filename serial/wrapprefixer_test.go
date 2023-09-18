package serial

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/adamcolton/luce/util/upgrade"
	"github.com/stretchr/testify/assert"
)

type stringPrefixer struct{}

func (stringPrefixer) PrefixReflectType(t reflect.Type, b []byte) ([]byte, error) {
	b = append(b, []byte(t.String())...)
	return b, nil
}

func (stringPrefixer) String() string {
	return "stringPrefixer"
}

func TestWrapPrefixer(t *testing.T) {
	r := strings.NewReplacer()
	wp := WrapPrefixer(stringPrefixer{})

	b, err := wp.PrefixInterfaceType(r, nil)
	assert.NoError(t, err)
	assert.Equal(t, "*strings.Replacer", string(b))

	// confirm we don't re-wrap a TypePrefixer
	_ = wp.(wrapPrefixer).ReflectTypePrefixer.(stringPrefixer)
	wp = WrapPrefixer(wp)
	_ = wp.(wrapPrefixer).ReflectTypePrefixer.(stringPrefixer)

	stringer, ok := upgrade.To[fmt.Stringer](wp)
	assert.True(t, ok)
	assert.Equal(t, stringPrefixer{}.String(), stringer.String())
}
