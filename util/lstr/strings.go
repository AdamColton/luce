package lstr

import (
	"strconv"
	"strings"

	"github.com/adamcolton/luce/util/liter"
)

var NumericReplacer = NewRemover(",", "$", "%")

// Strings is helpful when processing a list of strings, often the result of
// splitting. Fulfills liter.Iter.
type Strings struct {
	Strings         []string
	Err             error
	Preprocess      func(string) (skip bool, cleaned string)
	NumericReplacer *strings.Replacer

	idx int
	cur string
}

// DefaultPreprocess for Strings. Performs TrimSpace and skps if the string
// is empty.
var DefaultPreprocess = func(str string) (skip bool, cleaned string) {
	cleaned = strings.TrimSpace(str)
	skip = cleaned == ""
	return
}

// NewStrings from the provided strings with DefaultPreprocess.
func NewStrings(strs []string) *Strings {
	return (&Strings{
		Strings:         strs,
		Preprocess:      DefaultPreprocess,
		NumericReplacer: NumericReplacer,
	}).init()
}

// Len fulfills slice.Lener, returns length of Strings.
func (s *Strings) Len() int {
	return len(s.Strings)
}

// Next fulfill liter.Iter. Returns the next string.
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

// Cur fulfill liter.Iter, returns current string.
func (s *Strings) Cur() (str string, done bool) {
	return s.cur, s.Done()
}

// Done fulfills liter.Iter, returns true when iteration is done.
func (s *Strings) Done() bool {
	return s == nil || s.Err != nil || s.idx >= s.Len()
}

// Idx fulfills liter.Iter, returns current index.
func (s *Strings) Idx() int {
	return s.idx
}

// Start fulfills liter.Starter, resets the iteration.
func (s *Strings) Start() (str string, done bool) {
	s.idx = 0
	s.init()
	return s.Cur()
}

// Sub takes the current value of Strings, splits using the provided argument
// and creates a new instance of Strings. This can, for instance, be useful to
// have one instance of Strings that splits on newline and use Sub to split on
// commas.
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

// Float64 attempts to parse the current value as a float64. If it fails, the
// error is written to s.Err.
func (s *Strings) Float64() (f float64) {
	if s.Done() {
		return
	}
	str := s.NumericReplacer.Replace(liter.Pop(s))
	f, s.Err = strconv.ParseFloat(str, 64)
	return
}

// Int attempts to parse the current value as an int. If it fails, the
// error is written to s.Err.
func (s *Strings) Int() (i int) {
	if s.Done() {
		return
	}
	str := s.NumericReplacer.Replace(liter.Pop(s))
	i, s.Err = strconv.Atoi(str)
	return
}
