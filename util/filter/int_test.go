package filter

import (
	"testing"

	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestIntSlice(t *testing.T) {
	got := GTE(5).Slice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	expected := []int{5, 6, 7, 8, 9, 10}
	assert.Equal(t, expected, got)
}

func TestIntChan(t *testing.T) {
	ch := make(chan int)
	go func() {
		for _, i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
			ch <- i
		}
		close(ch)
	}()

	to := timeout.After(5, func() {
		expected := []int{5, 6, 7, 8, 9, 10}
		get := GTE(5).Chan(ch, 0)
		for _, e := range expected {
			assert.Equal(t, e, <-get)
		}
	})
	assert.NoError(t, to)
}

func TestIntBools(t *testing.T) {
	tt := map[string]struct {
		f Filter[int]
		x map[int]bool
	}{
		"4<x_AND_x<7": {
			f: LT(7).And(GT(4)),
			x: map[int]bool{
				4: false,
				5: true,
				6: true,
				7: false,
			},
		},
		"4>x_OR_x>7": {
			f: GT(7).Or(LT(4)),
			x: map[int]bool{
				4: false,
				3: true,
				8: true,
				7: false,
			},
		},
		"!(x>5)": {
			f: GT(5).Not(),
			x: map[int]bool{
				5: true,
				6: false,
				7: false,
				4: true,
				3: true,
			},
		},
		"x==3": {
			f: EQ(3),
			x: map[int]bool{
				2: false,
				3: true,
				4: false,
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
