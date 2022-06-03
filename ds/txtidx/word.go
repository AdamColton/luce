package txtidx

import (
	"strings"

	"github.com/adamcolton/luce/ds/bintrie"
)

type word struct {
	str string
	wordIDX
	Documents *docSet
}

type wordIDX uint32

// str must start with letterNumber but can have trailing non-letter number
func root(str string) string {
	s := newScanner(str)
	s.matchLetterNumber(false)
	return strings.ToLower(s.str(0, s.idx))
}

type words []*word

func (ws words) docSetUnion() *docSet {
	ln := len(ws)
	if ln == 0 {
		return newDocSet()
	} else if ln == 1 {
		return ws[0].Documents
	}
	out := ws[0].Documents.union(ws[1].Documents)
	for _, w := range ws[2:] {
		out.merge(w.Documents)
	}
	return out
}

func (ws words) bintrieMulti() bintrie.Multi {
	ln := len(ws)
	if ln == 0 {
		return nil
	}
	out := make(bintrie.Multi, ln)
	for i, w := range ws {
		out[i] = w.Documents.t
	}
	return out
}

func (ws words) largestDocSet() int {
	var max int
	for _, w := range ws {
		if s := w.Documents.t.Size(); s > max {
			max = s
		}
	}
	return max
}
