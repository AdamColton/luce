package fngraph

import (
	"image"
	"image/color"

	"github.com/adamcolton/luce/math"
	"github.com/adamcolton/luce/math/cmpr"
	"github.com/adamcolton/luce/math/funcs"
	"github.com/adamcolton/luce/math/ints"
	"github.com/adamcolton/luce/math/numiter"
)

// Channel provides a color channel, x and y will always be between 0 and 1.
type Channel func(x, y float64) float64

type Range func(x float64) float64

func NewRange(start, end float64) Range {
	d := end - start
	return func(x float64) float64 {
		return x*d + start
	}
}

type MapInput func(x, y float64) []float64

func NewMapInput(s []float64, xIdx, yIdx int, xRng, yRng Range) MapInput {
	return func(x, y float64) []float64 {
		s[xIdx] = xRng(x)
		s[yIdx] = yRng(y)
		return s
	}
}

func (mi MapInput) Channel(m funcs.M) Channel {
	return func(x, y float64) float64 {
		return m(mi(x, y))
	}
}

const (
	R = iota
	G
	B
	All
)

type Graph struct {
	Channels [3]Channel
}

func New(r, g, b Channel) Graph {
	return Graph{
		Channels: [3]Channel{r, g, b},
	}
}

type Matrix [][]float64

type ColorMatrix [3]Matrix

func (g Graph) ColorMatrix(w, h int) ColorMatrix {
	var m ColorMatrix
	for i := range All {
		if c := g.Channels[i]; c != nil {
			m[i] = c.Matrix(w, h).Normalize()
		}
	}
	return m
}

const (
	max16  = ints.MaxU16
	max16f = float64(ints.MaxU16)
)

var upLeft = image.Point{0, 0}

func (g Graph) Image(w, h int) image.Image {
	cm := g.ColorMatrix(w, h)
	lowRight := image.Point{w, h}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	buf := [3]uint16{0, 0, 0}
	numiter.IntGrid(w, h).Iter().For(func(t []int) {
		x, y := t[0], t[1]
		for i := range buf {
			if cm[i] == nil {
				buf[i] = 0
			} else {
				buf[i] = uint16(max16f * cm[i][x][y])
			}
		}
		c := color.RGBA64{buf[R], buf[G], buf[B], max16}
		img.Set(x, y, c)
	})
	return img
}

func (c Channel) Matrix(w, h int) Matrix {
	xr := numiter.Steps(0.0, 1.0, uint(w), true).Iter()
	yr := numiter.Steps(0.0, 1.0, uint(h), true).Iter()

	out := make(Matrix, w)
	for xi, xf := range xr.Seq2 {
		out[xi] = make([]float64, h)
		for yi, yf := range yr.Seq2 {
			out[xi][yi] = c(xf, yf)
		}
	}
	return out
}

func (m Matrix) MinMax() (min, max float64) {
	min, max = m[0][0], m[0][0]

	for _, col := range m {
		max = cmpr.Max(max, cmpr.MaxN(col...))
		min = cmpr.Max(min, cmpr.MinN(col...))
	}

	return
}

func (m Matrix) Apply(fn func(float64) float64, buf Matrix) Matrix {
	if buf == nil {
		buf = m
	}
	for x, col := range m {
		for y, v := range col {
			buf[x][y] = fn(v)
		}
	}
	return buf
}

func (m Matrix) Normalize() Matrix {
	min, max := m.MinMax()
	scale := math.Scale(min, max)
	m.Apply(scale, m)
	return m
}
