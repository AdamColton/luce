package lstr

import (
	"regexp"

	"github.com/adamcolton/luce/ds/morph"
	"github.com/adamcolton/luce/util/liter"
)

// RegexpSlice is a list of Regexps that are tried in order.
type RegexpSlice []*regexp.Regexp

var mustRe = morph.NewValAll(regexp.MustCompile)

// NewRegexpSlice compiles each string with regexp.MustCompile, so it panics if
// one is not a valid expression.
func NewRegexpSlice(strs ...string) RegexpSlice {
	return RegexpSlice(mustRe.Slice(strs, nil))
}

// Match returns the string of the first Regexp that matches str and the
// submatches (see regexp.FindStringSubmatch). If none matches it returns "" and
// nil.
func (nrs RegexpSlice) Match(str string) (string, []string) {
	for _, re := range nrs {
		m := re.FindStringSubmatch(str)
		if len(m) > 0 {
			return re.String(), m
		}
	}
	return "", nil
}

// MatchIter calls Match on the values of the iterator, starting with the
// current one, until one matches, and returns that result. The iterator is left
// on the value after the match. If the iterator is done or nothing matches, it
// returns "" and nil.
func (nrs RegexpSlice) MatchIter(it liter.Iter[string]) (id string, m []string) {
	if it.Done() {
		return "", nil
	}
	liter.Seek(it, func(str string) (done bool) {
		if id != "" {
			return true
		}
		id, m = nrs.Match(str)
		return false
	})
	return
}
