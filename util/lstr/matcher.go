package lstr

import (
	"unicode"
)

// Matcher returns true when a rune matches some criteria.
type Matcher interface {
	Matches(r rune) bool
}

// Rune fulfills Matcher returning true for an exact match
type Rune rune

// Matches fulfills Matcher returning true when Rune == r
func (rn Rune) Matches(r rune) bool {
	return rune(rn) == r
}

// Range fulfills Matcher returning true when when given a rune that is between
// the two values (inclusivly).
type Range [2]rune

// Matches fulfills Matcher returning true when when given a rune that is between
// the two values (inclusivly).
func (rng Range) Matches(r rune) bool {
	return r >= rng[0] && r <= rng[1]
}

// Not takes a Matcher and inverts the output.
type Not struct {
	Matcher
}

// Matches inverts the output of the underlying Matcher.
func (n Not) Matches(r rune) bool {
	return !n.Matcher.Matches(r)
}

// Or takes many Matcher and returns true if any of them return true for a rune.
type Or []Matcher

// Matches returns true if any of the Matchers return true for r.
func (o Or) Matches(r rune) bool {
	for _, m := range o {
		if m.Matches(r) {
			return true
		}
	}
	return false
}

// And takes many Matcher and returns true if all of them return true for a
// rune.
type And []Matcher

// Matches returns true if all of the Matchers return true for a rune.
func (a And) Matches(r rune) bool {
	for _, m := range a {
		if !m.Matches(r) {
			return false
		}
	}
	return true
}

// MatcherFunc wraps a func so it fulfills Matcher.
type MatcherFunc func(r rune) bool

// Matches calls the underlying MatcherFunc.
func (fn MatcherFunc) Matches(r rune) bool {
	return fn(r)
}

var (
	IsLetter  MatcherFunc = unicode.IsLetter
	IsNumber  MatcherFunc = unicode.IsNumber
	IsDigit   MatcherFunc = unicode.IsDigit
	IsControl MatcherFunc = unicode.IsControl
	IsGraphic MatcherFunc = unicode.IsGraphic
	IsLower   MatcherFunc = unicode.IsLower
	IsMark    MatcherFunc = unicode.IsMark
	IsPrint   MatcherFunc = unicode.IsPrint
	IsPunct   MatcherFunc = unicode.IsPunct
	IsSpace   MatcherFunc = unicode.IsSpace
	IsTitle   MatcherFunc = unicode.IsTitle
	IsUpper   MatcherFunc = unicode.IsUpper
)
