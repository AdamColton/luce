package lstr

import "strings"

// Strings is helpful when processing a list of strings, often the result of
// splitting. Fulfills liter.Iter and slice.Lener.
type Strings struct {
	// Strings are the strings being iterated over.
	Strings []string
	// Err is the first error from parsing a value. Once it is set, the Strings
	// are done.
	Err error
	// Preprocess is applied to each string. If skip is true the string is
	// skipped, otherwise cleaned is the value. It can be nil.
	Preprocess func(string) (skip bool, cleaned string)

	idx int
	cur string
}

// DefaultPreprocess for Strings. Performs TrimSpace and skips if the string
// is empty.
var DefaultPreprocess = func(str string) (skip bool, cleaned string) {
	cleaned = strings.TrimSpace(str)
	skip = cleaned == ""
	return
}

// NewStrings from the provided strings with DefaultPreprocess.
func NewStrings(strs []string) *Strings {
	return (&Strings{
		Strings:    strs,
		Preprocess: DefaultPreprocess,
	}).init()
}

// Len fulfills slice.Lener, returns the number of strings, including any that
// Preprocess skips.
func (s *Strings) Len() int {
	return len(s.Strings)
}

// Next fulfills liter.Iter. It moves to the next string that Preprocess does
// not skip and returns it and if the iteration is done.
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

// Cur fulfills liter.Iter, returns current string.
func (s *Strings) Cur() (str string, done bool) {
	return s.cur, s.Done()
}

// Done fulfills liter.Iter, returns true when iteration is done, which includes
// when Err is set or when the Strings is nil.
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
