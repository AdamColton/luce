package lfile

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiglob(t *testing.T) {
	restore := Glob
	defer func() { Glob = restore }()

	mockdir := map[string][]string{
		"foo*":  {"foo", "fooer", "foo.bar"},
		"*.bar": {"foo.bar", "bar.bar"},
	}
	Glob = func(pattern string) ([]string, error) {
		matches, found := mockdir[pattern]
		if found {
			return matches, nil
		}
		return nil, filepath.ErrBadPattern
	}

	got, err := MultiGlob{"foo*", "*.bar"}.Filenames()
	assert.NoError(t, err)
	sort.Strings(got)
	expected := []string{"bar.bar", "foo", "foo.bar", "fooer"}
	assert.Equal(t, expected, got)

	_, err = MultiGlob{"]]] == bad pattern == [[["}.Filenames()
	assert.Equal(t, filepath.ErrBadPattern, err)
}
