package lstr_test

import (
	"testing"

	"github.com/adamcolton/luce/util/lstr"
	"github.com/stretchr/testify/assert"
)

func TestCamelCase(t *testing.T) {
	tt := map[string][]string{
		"AbcdEfgh": {"Abcd", "Efgh"},
		"AbcdEFgh": {"Abcd", "E", "Fgh"},
		"AbcDEFgh": {"Abc", "DE", "Fgh"},
		"aBcDeFgh": {"a", "Bc", "De", "Fgh"},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc, lstr.CamelCase(n).Slice(n))
		})
	}
}
