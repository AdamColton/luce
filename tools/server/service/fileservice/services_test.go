package fileservice

import (
	"testing"

	"github.com/adamcolton/luce/util/lfile/lfilemock"
)

func TestServices(t *testing.T) {
	repo = lfilemock.Parse(map[string]any{
		"dir1": map[string]any{},
		"dir2": map[string]any{},
	})
	srvs := New()
	srvs.New("foo")
}
