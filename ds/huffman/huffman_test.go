package huffman

import (
	"fmt"
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/serial/rye"
	"github.com/stretchr/testify/assert"
)

var letters = lmap.Map[rune, int]{
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

	{
		lt := slice.LT[rune]()
		expected := lt.Sort(letters.Keys(nil))
		got := make([]rune, 0, len(expected))
		ht.All(func(r rune) {
			got = append(got, r)
		})
		lt.Sort(got)
		assert.Equal(t, expected, got)
	}

	assert.Equal(t, len(data), ht.Len())

	l := NewLookup(ht)
	expectedBits := &rye.Bits{
		Ln:   4,
		Data: []byte{7},
	}
	assert.Equal(t, expectedBits, l.Get('A'))
	expectedBits = &rye.Bits{
		Ln:   9,
		Data: []byte{52, 0},
	}
	assert.Equal(t, expectedBits, l.Get('Z'))

	// round trip
	expected := []rune("THISISATEST")

	enc := Encode[rune](list.Slice(expected), l)

	it, _, _ := ht.Iter(enc).Factory()
	got := slice.IterSlice[rune](it, nil)
	assert.Equal(t, expected, got)
	assert.True(t, enc.Ln < len(expected)*8)

	var less = slice.LT[rune]()
	expectedRunes := less.Sort(letters.Keys(nil))
	gotRunes := less.Sort(l.All())
	assert.Equal(t, expectedRunes, gotRunes)
}

func TestTranslate(t *testing.T) {
	data := make([]Frequency[[]byte], 0, len(letters))
	var expectedSlices [][]byte
	for r, c := range letters {
		b := []byte(string(r))
		data = append(data, Frequency[[]byte]{Val: b, Count: c})
		expectedSlices = append(expectedSlices, b)
	}
	ht := New(data)
	l := NewTranslateLookup(ht, func(b []byte) string {
		return string(b)
	})
	expectedBits := &rye.Bits{
		Ln:   4,
		Data: []byte{7},
	}
	assert.Equal(t, expectedBits, l.Get([]byte("A")))
	expectedBits = &rye.Bits{
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
	enc := Encode[[]byte](list.Slice(expected), l)

	got := slice.IterSlice[[]byte](ht.Iter(enc), nil)
	assert.Equal(t, expected, got)
	assert.True(t, enc.Ln < len(expected)*8)

	less := slice.Less[[]byte](func(i, j []byte) bool {
		return string(i) < string(j)
	})
	gotSlices := less.Sort(l.All())

	assert.Equal(t, less.Sort(expectedSlices), gotSlices)
}

func ExampleEncode_roundTrip() {
	var letterFrequencyMap = map[rune]int{
		'E': 21912, 'T': 16587, 'A': 14810, 'O': 14003, 'I': 13318, 'N': 12666,
		'S': 11450, 'R': 10977, 'H': 10795, 'D': 7874, 'L': 7253, 'U': 5246,
		'C': 4943, 'M': 4761, 'F': 4200, 'Y': 3853, 'W': 3819, 'G': 3693,
		'P': 3316, 'B': 2715, 'V': 2019, 'K': 1257, 'X': 315, 'Q': 205,
		'J': 188, 'Z': 128,
	}
	tree := MapNew(letterFrequencyMap)
	l := NewLookup(tree)

	bits := Encode[rune](list.Slice([]rune("THISISATEST")), l)
	// Encoded length is 5, much less than the 11 characters
	fmt.Println("Length:", len(bits.Data))

	str := string(slice.IterSlice[rune](tree.Iter(bits), nil))
	fmt.Println(str)

	// Output:
	// Length: 5
	// THISISATEST
}
