package ldate_test

import (
	"testing"

	"github.com/adamcolton/luce/util/ldate"
	"github.com/stretchr/testify/assert"
)

func TestFromStringErrors(t *testing.T) {
	for _, str := range []string{"", "2024", "2024_01"} {
		_, err := ldate.FromString(str)
		assert.Equal(t, ldate.ErrMalformed, err, str)
	}
	for _, str := range []string{"x_01_01", "2024_x_01", "2024_01_x"} {
		_, err := ldate.FromString(str)
		assert.Error(t, err, str)
		assert.NotEqual(t, ldate.ErrMalformed, err, str)
	}

	d, err := ldate.FromString("-001_12_31")
	assert.NoError(t, err)
	assert.Equal(t, ldate.New(-1, 12, 31), d)
}
