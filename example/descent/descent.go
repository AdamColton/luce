package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/math/funcs"
	"github.com/adamcolton/luce/math/ints/comb"
	"github.com/fogleman/gg"
)

var (
	size  = 500
	s64   = float64(size)
	sizeI = list.Generator[float64]{
		Length: size,
		Fn: func(i int) float64 {
			return float64(i) / s64
		},
	}
	grid = list.Combinator(sizeI, sizeI, comb.Cross[int])

	s2    = s64 / 2
	scale = func(x float64) float64 {
		return x*100 + s2
	}
	circA = func(t float64) (x, y float64) {
		x, y = math.Sincos(t * 6.28)
		return scale(x), scale(y)
	}
	circB = func(t float64) (x, y float64) {
		x, y = math.Sincos(t*6.28 + 1)
		x += 0.5
		return scale(x), scale(y)
	}
	dist funcs.M = func(t []float64) float64 {
		x0, y0 := circA(t[0])
		x1, y1 := circB(t[1])
		dx := (x0 - x1)
		dy := (y0 - y1)
		d := math.Sqrt(dx*dx + dy*dy)
		return d
	}
	resolution = 100.0
	steps      = list.Generator[float64]{
		Length: int(resolution),
		Fn: func(i int) float64 {
			return float64(i) / resolution
		},
	}.Wrap()
)

func clr(ctx *gg.Context) {
	ctx.SetColor(color.RGBA{0, 0, 0, 255})
	ctx.Clear()
}

func main() {
	ctx := gg.NewContext(size, size)
	clr(ctx)

	s := func(x float64) int {
		return int(s64 * x)
	}

	grid.Iter().For(func(t struct {
		A float64
		B float64
	}) {
		x, y := t.A, t.B

		sx, sy := s(x), s(y)
		g := dist([]float64{x, y})
		ctx.SetColor(color.RGBA{0, uint8(g), 0, 255})
		ctx.SetPixel(sx, sy)
	})

	s2 := func(x float64) float64 {
		return x * s64
	}

	// x^n = .1
	// .1^1/.n

	rand.Seed(time.Now().UnixMicro())
	d := (&funcs.Descender{
		Multi: funcs.Multi{
			Ln: 2,
			M:  dist,
		},
		Steps: 100,
		X:     []float64{rand.Float64(), rand.Float64()},
	}).Init()

	fmt.Println(d.DG)
	ctx.MoveTo(s2(d.X[0]), s2(d.X[1]))

	for d.Steps > 0 {
		d.Step()
		ctx.SetColor(color.RGBA{255 - uint8(d.Steps), 0, uint8(d.Steps), 255})
		x, y := s2(d.X[0]), s2(d.X[1])
		ctx.LineTo(x, y)
		ctx.Stroke()
		ctx.MoveTo(x, y)

		//d.G = cmpr.Min(d.G, setG(0.1, d.DX))
	}

	ctx.SavePNG("steps.png")

	clr(ctx)
	ctx.SetColor(color.RGBA{255, 0, 0, 255})
	draw(circA, ctx)
	ctx.SetColor(color.RGBA{0, 0, 255, 255})
	draw(circB, ctx)

	xa, ya := circA(d.X[0])
	ctx.SetColor(color.RGBA{255, 255, 0, 255})
	ctx.DrawCircle(xa, ya, 3)
	ctx.Stroke()

	xb, yb := circB(d.X[1])
	ctx.SetColor(color.RGBA{0, 255, 255, 255})
	ctx.DrawCircle(xb, yb, 3)
	ctx.Stroke()

	ctx.SavePNG("result.png")

}

func draw(fn func(float64) (float64, float64), ctx *gg.Context) {
	i := steps.Iter()
	x0, y0 := fn(i.Pop())
	ctx.MoveTo(x0, y0)
	i.For(func(t float64) {
		ctx.LineTo(fn(t))
	})
	ctx.LineTo(x0, y0)
	ctx.Stroke()
}
