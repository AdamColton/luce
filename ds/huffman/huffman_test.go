package huffman

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var letters = map[rune]int{
	'E': 21912,
	'T': 16587,
	'A': 14810,
	'O': 14003,
	'I': 13318,
	'N': 12666,
	'S': 11450,
	'R': 10977,
	'H': 10795,
	'D': 7874,
	'L': 7253,
	'U': 5246,
	'C': 4943,
	'M': 4761,
	'F': 4200,
	'Y': 3853,
	'W': 3819,
	'G': 3693,
	'P': 3316,
	'B': 2715,
	'V': 2019,
	'K': 1257,
	'X': 315,
	'Q': 205,
	'J': 188,
	'Z': 128,
}

func TestFromMap(t *testing.T) {
	data := letters
	ht := MapNew(data)
	assert.NotNil(t, ht)

	l := NewLookup(ht)
	expectedBits := &Bits{
		Ln:   4,
		Data: []byte{7},
	}
	assert.Equal(t, expectedBits, l.Get('A'))
	expectedBits = &Bits{
		Ln:   9,
		Data: []byte{52, 0},
	}
	assert.Equal(t, expectedBits, l.Get('Z'))

	// round trip
	expected := []rune("THISISATEST")
	enc := Encode(expected, l)
	got := ht.ReadAll(enc)
	assert.Equal(t, expected, got)
	assert.True(t, enc.Ln < len(expected)*8)
}

func TestTranslate(t *testing.T) {
	data := make([]Frequency[[]byte], 0, len(letters))
	for r, c := range letters {
		b := []byte(string(r))
		data = append(data, Frequency[[]byte]{Val: b, Count: c})
	}
	ht := New(data)
	l := NewTranslateLookup(ht, func(b []byte) string {
		return string(b)
	})
	expectedBits := &Bits{
		Ln:   4,
		Data: []byte{7},
	}
	assert.Equal(t, expectedBits, l.Get([]byte("A")))
	expectedBits = &Bits{
		Ln:   9,
		Data: []byte{52, 0},
	}
	assert.Equal(t, expectedBits, l.Get([]byte("Z")))

	// round trip
	expected := [][]byte{
		[]byte("T"),
		[]byte("H"),
		[]byte("I"),
		[]byte("S"),
		[]byte("I"),
		[]byte("S"),
		[]byte("A"),
		[]byte("T"),
		[]byte("E"),
		[]byte("S"),
		[]byte("T"),
	}
	enc := Encode(expected, l)
	got := ht.ReadAll(enc)
	assert.Equal(t, expected, got)
	assert.True(t, enc.Ln < len(expected)*8)
}
