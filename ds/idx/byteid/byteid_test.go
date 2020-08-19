package byteid_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid"
	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	l := byteid.IDLen(5)
	assert.Equal(t, byteid.ID{0, 0, 0, 0, 0}, l.Zero())
	assert.Len(t, l.Rand(), int(l))
}

func TestIncDec(t *testing.T) {
	tt := map[string]struct {
		in       byteid.ID
		expected byteid.ID
	}{
		"basic": {
			in:       byteid.ID{1, 2, 3},
			expected: byteid.ID{1, 2, 4},
		},
		"wrap-once": {
			in:       byteid.ID{1, 2, 255},
			expected: byteid.ID{1, 3, 0},
		},
		"wrap-all": {
			in:       byteid.ID{255, 255, 255, 255},
			expected: byteid.ID{0, 0, 0, 0},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.in.Inc())
			assert.Equal(t, tc.in, tc.expected.Dec())
		})
	}
}

func TestFuzzIncDec(t *testing.T) {
	l := byteid.IDLen(4)
	for i := 0; i < 1000; i++ {
		id := l.Rand()
		assert.Equal(t, id, id.Inc().Dec())
		assert.Equal(t, id, id.Dec().Inc())
	}
}
