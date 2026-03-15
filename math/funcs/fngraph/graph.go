package fngraph

import (
	"image"
	"image/color"

	"github.com/adamcolton/luce/math/ints"
	"github.com/adamcolton/luce/math/numiter"
)

type Graph struct {
	RGBChannels [3]Channel
}

func New(r, g, b Channel) Graph {
	return Graph{
		RGBChannels: [3]Channel{r, g, b},
	}
}

type ColorMatrix [3]Matrix

func (g Graph) ColorMatrix(w, h int) ColorMatrix {
	var m ColorMatrix
	for i := range All {
		if c := g.RGBChannels[i]; c != nil {
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
