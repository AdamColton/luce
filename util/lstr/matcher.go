package lstr

import "unicode"

type Matcher interface {
	Matches(r rune) bool
}

type Range [2]rune

func (rng Range) Matches(r rune) bool {
	return r >= rng[0] && r <= rng[1]
}

type Not struct {
	Matcher
}

func (n Not) Matches(r rune) bool {
	return !n.Matcher.Matches(r)
}

type Or []Matcher

func (o Or) Matches(r rune) bool {
	for _, m := range o {
		if m.Matches(r) {
			return true
		}
	}
	return false
}

type And []Matcher

func (a And) Matches(r rune) bool {
	for _, m := range a {
		if !m.Matches(r) {
			return false
		}
	}
	return true
}

type MatcherFunc func(r rune) bool

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
