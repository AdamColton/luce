package scratchidx_test

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store/entity"
	"github.com/adamcolton/luce/store/entity/scratchidx"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/stretchr/testify/assert"
)

type Foo struct {
	id   byte
	Word string
}

func (f *Foo) EntKey() []byte {
	return []byte{f.id}
}

func TestIdxr(t *testing.T) {
	cmpr := func(a, b string) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}
	getKey := func(f *Foo) string {
		return f.Word
	}
	idx := scratchidx.NewIndex(getKey, cmpr, nil)

	s := lerr.Must(ephemeral.Factory(bytebtree.New, 1).Store([]byte("test")))
	b := entity.NewGobBuilder(s)
	lerr.Must(idx.Store(b))

	words := []string{
		"young", "unity", "defend", "storage", "law", "pack", "strike",
		"triangle", "agenda", "knee", "model", "resist", "hike", "aspect",
		"wander", "photography", "strain", "school", "definite", "advocate",
		"map", "projection", "warm", "research", "instinct", "parking",
		"contain", "danger", "deadly", "premature", "day", "brilliance",
		"diplomatic", "colony", "effort", "faith", "harbor", "weigh",
		"impound", "bond", "acquit", "apparatus", "tile", "heart", "wait",
	}

	for i, w := range words {
		e := &Foo{
			id:   byte(i),
			Word: w,
		}
		idx.Update(e)
	}

	got := idx.Get("warm")
	expected := [][]byte{{22}}
	assert.Equal(t, expected, got)

	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(idx)
	assert.NoError(t, err)

	idx2 := scratchidx.NewIndex(getKey, cmpr, nil)
	buf = bytes.NewBuffer(buf.Bytes())
	lerr.Panic(gob.NewDecoder(buf).Decode(idx2))

	// When idx gets encoded, root is a graph.Ptr and it's encoded by reference
	// so when it's decoded, it's nil.
	got = idx2.Get("warm")
	assert.Equal(t, expected, got)
}
