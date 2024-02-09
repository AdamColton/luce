package lstr

import (
	"strconv"
	"strings"

	"github.com/adamcolton/luce/util/liter"
)

var NumericReplacer = NewRemover(",", "$", "%")

// Strings is helpful when processing a list of strings, often the result of
// splitting.
type Strings struct {
	Strings         []string
	Err             error
	Preprocess      func(string) (skip bool, cleaned string)
	NumericReplacer *strings.Replacer

	idx int
	cur string
}

var DefaultPreprocess = func(str string) (skip bool, cleaned string) {
	cleaned = strings.TrimSpace(str)
	skip = cleaned == ""
	return
}

func NewStrings(strs []string) *Strings {
	return (&Strings{
		Strings:         strs,
		Preprocess:      DefaultPreprocess,
		NumericReplacer: NumericReplacer,
	}).init()
}

func (s *Strings) Len() int {
	return len(s.Strings)
}

func (s *Strings) Next() (str string, done bool) {
	done = s.inc()
	for ; !done && s.setCur(); done = s.inc() {
	}
	if done {
		s.cur = ""
	}

	return s.cur, done
}

func (s *Strings) inc() (done bool) {
	s.idx++
	return s.Done()
}

func (s *Strings) setCur() (skip bool) {
	s.cur = s.Strings[s.idx]
	if s.Preprocess != nil {
		skip, s.cur = s.Preprocess(s.cur)
	}
	return skip
}

func (s *Strings) init() *Strings {
	if !s.Done() && s.setCur() {
		s.Next()
	}
	return s
}

func (s *Strings) Cur() (str string, done bool) {
	return s.cur, s.Done()
}

func (s *Strings) Done() bool {
	return s == nil || s.Err != nil || s.idx >= s.Len()
}

func (s *Strings) Idx() int {
	return s.idx
}

func (s *Strings) Start() (str string, done bool) {
	s.idx = 0
	s.init()
	return s.Cur()
}

func (s *Strings) Sub(split string) *Strings {
	done := s.Done()
	if done {
		return nil
	}
	strs := strings.Split(liter.Pop(s), split)
	sub := (&Strings{
		Strings:         strs,
		Preprocess:      s.Preprocess,
		NumericReplacer: s.NumericReplacer,
	}).init()
	if sub.Done() {
		return s.Sub(split)
	}
	return sub
}

func (s *Strings) Float64() (f float64) {
	if s.Done() {
		return
	}
	str := s.NumericReplacer.Replace(liter.Pop(s))
	f, s.Err = strconv.ParseFloat(str, 64)
	return
}

func (s *Strings) Int() (i int) {
	if s.Done() {
		return
	}
	str := s.NumericReplacer.Replace(liter.Pop(s))
	i, s.Err = strconv.Atoi(str)
	return
}
