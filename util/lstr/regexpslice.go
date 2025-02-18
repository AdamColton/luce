package lstr

import (
	"regexp"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/liter"
)

type RegexpSlice []*regexp.Regexp

var mustRe = slice.ForAll(regexp.MustCompile)

func NewRegexpSlice(strs ...string) RegexpSlice {
	return RegexpSlice(mustRe.Slice(strs, nil))
}

func (nrs RegexpSlice) Match(str string) (string, []string) {
	for _, re := range nrs {
		m := re.FindStringSubmatch(str)
		if len(m) > 0 {
			return re.String(), m
		}
	}
	return "", nil
}

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
