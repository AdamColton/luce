package comb_test

import (
	"testing"

	"github.com/adamcolton/luce/math/comb"
	"github.com/stretchr/testify/assert"
)

// Classic Binomial function, no optimizations. Use for benchmark and validating
// memo values.
func Classic(n, i int) int {
	// https://math.stackexchange.com/questions/202554/how-do-i-compute-binomial-coefficients-efficiently
	if n < 0 || i > n || i < 0 {
		return 0
	} else if i == 0 {
		return 1
	} else if i > n/2 {
		return Classic(n, n-i)
	}

	return n * Classic(n-1, i-1) / i
}

type key [2]int

var mapMemo = make(map[key]int)

// MapMemo is an earlier version that got moved into the tests as a benchmark
// comparison.
func MapMemo(n, i int) int {
	if n < 0 || i > n || i < 0 {
		return 0
	}
	if i > n/2 {
		i = n - i
	}
	if i < 2 {
		return (i * n) + (1 - i)
	}

	k := key{n, i}
	if v, found := mapMemo[k]; found {
		return v
	}

	v := n * MapMemo(n-1, i-1) / i
	mapMemo[k] = v
	return v
}

func TestAgainstClassic(t *testing.T) {
	old := comb.Memo
	comb.Memo = make([]int, 128)
	assert.Equal(t, Classic(35, 17), comb.Binomial(35, 17))

	// Two passes, once to refill the memo, once to test with the memo full
	for passes := 0; passes < 2; passes++ {
		for n := 0; n < 35; n++ {
			for i := 0; i <= n; i++ {
				b := comb.Binomial(n, i)
				k := Classic(n, i)
				if b != k {
					t.Error(n, i)
				}
			}
		}
	}
	assert.Equal(t, old, comb.Memo[:256])

}

func TestOutOfBounds(t *testing.T) {
	assert.Equal(t, 0, comb.Binomial(-1, 1))
}

var maxN = 256

func BenchmarkMemo(b *testing.B) {
	n := 0
	i := -1
	for iter := 0; iter < b.N; iter++ {
		i++
		if i > n {
			n, i = n+1, 0
			if n > maxN {
				n, i = 0, 0
			}
		}
		comb.Binomial(n, i)
	}
}

func BenchmarkMap(b *testing.B) {
	n := 0
	i := -1
	for iter := 0; iter < b.N; iter++ {
		i++
		if i > n {
			n, i = n+1, 0
			if n > maxN {
				n, i = 0, 0
			}
		}
		MapMemo(n, i)
	}
}

func BenchmarkClassic(b *testing.B) {
	n := 0
	i := -1
	for iter := 0; iter < b.N; iter++ {
		i++
		if i > n {
			n, i = n+1, 0
			if n > maxN {
				n, i = 0, 0
			}
		}
		Classic(n, i)
	}
}
