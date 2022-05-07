package filter

import (
	"testing"

	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestFloatSlice(t *testing.T) {
	got := GTE(5.0).Slice([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	expected := []float64{5, 6, 7, 8, 9, 10}
	assert.Equal(t, expected, got)
}

func TestFloatChan(t *testing.T) {
	ch := make(chan float64)
	go func() {
		for _, i := range []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
			ch <- i
		}
		close(ch)
	}()

	to := timeout.After(5, func() {
		expected := []float64{5, 6, 7, 8, 9, 10}
		get := GTE(5.0).Chan(ch, 0)
		for _, e := range expected {
			assert.Equal(t, e, <-get)
		}
	})
	assert.NoError(t, to)
}

func TestFloatBools(t *testing.T) {
	tt := map[string]struct {
		f Filter[float64]
		x map[float64]bool
	}{
		"4<x_AND_x<7": {
			f: LT(7.0).And(GT(4.0)),
			x: map[float64]bool{
				4: false,
				5: true,
				6: true,
				7: false,
			},
		},
		"4>x_OR_x>7": {
			f: GT(7.0).Or(LT(4.0)),
			x: map[float64]bool{
				4: false,
				3: true,
				8: true,
				7: false,
			},
		},
		"!(x>5)": {
			f: GT(5.0).Not(),
			x: map[float64]bool{
				5: true,
				6: false,
				7: false,
				4: true,
				3: true,
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
