package filter

import (
	"testing"

	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestLazyOr(t *testing.T) {
	s := Filter[string](func(s string) bool { return true }).
		Or(func(s string) bool {
			t.Error("you should not be here")
			return false
		})
	s("test")
}

func TestStringSlice(t *testing.T) {
	got := GTE("5").Slice([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"})
	expected := []string{"5", "6", "7", "8", "9"}
	assert.Equal(t, expected, got)
}

func TestStringChan(t *testing.T) {
	ch := make(chan string)
	go func() {
		for _, i := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"} {
			ch <- i
		}
		close(ch)
	}()

	to := timeout.After(5, func() {
		expected := []string{"5", "6", "7", "8", "9"}
		get := GTE("5").Chan(ch, 0)
		for _, e := range expected {
			assert.Equal(t, e, <-get)
		}
	})
	assert.NoError(t, to)
}

func TestStringBools(t *testing.T) {
	tt := map[string]struct {
		f Filter[string]
		x map[string]bool
	}{
		"4<x_AND_x<7": {
			f: LT("7").And(GT("4")),
			x: map[string]bool{
				"4": false,
				"5": true,
				"6": true,
				"7": false,
			},
		},
		"4>x_OR_x>7": {
			f: GT("7").Or(LT("4")),
			x: map[string]bool{
				"4": false,
				"3": true,
				"8": true,
				"7": false,
			},
		},
		"!(x>5)": {
			f: GT("5").Not(),
			x: map[string]bool{
				"5": true,
				"6": false,
				"7": false,
				"4": true,
				"3": true,
			},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			for i, b := range tc.x {
				assert.Equal(t, b, tc.f(i))
			}
		})
	}
}

func TestHelpers(t *testing.T) {
	f := Prefix("test")
	assert.True(t, f("testing"))
	assert.False(t, f("nottesting"))

	f = Suffix("test")
	assert.True(t, f("atest"))
	assert.False(t, f("testnot"))

	f = Contains("test")
	assert.True(t, f("atests"))
	assert.False(t, f("itdoesnot"))
}

func TestRegex(t *testing.T) {
	tt := map[string]map[string]bool{
		"ca*t": {
			"cat":         true,
			"ct":          true,
			"cot":         false,
			"acat":        true,
			"dogcatmouse": true,
		},
		"^ca*t$": {
			"cat":         true,
			"ct":          true,
			"cot":         false,
			"acat":        false,
			"dogcatmouse": false,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			r := MustRegex(n)
			for s, expected := range tc {
				assert.Equal(t, expected, r(s))
			}
			var err error
			r, err = Regex(n)
			assert.NoError(t, err)
			for s, expected := range tc {
				assert.Equal(t, expected, r(s))
			}
		})
	}

	_, err := Regex("bad [ regex")
	assert.Error(t, err)
}
