// Command descent draws how gradient descent, from math/funcs, looks for the
// closest points of two circles: steps.png is the path it takes over the map of
// distances and result.png shows the circles and the two points it found.
package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/adamcolton/luce/math/funcs"
	"github.com/adamcolton/luce/math/numiter"
	"github.com/fogleman/gg"
)

var (
	size = 500
	s64  = float64(size)
	grid = numiter.Grid(0.0, 1.0, 1/s64, 0.0, 1.0, 1/s64)

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
	steps      = numiter.NewRange(0, 1, 1/resolution)
)

// clr clears the context to black.
func clr(ctx *gg.Context) {
	ctx.SetColor(color.RGBA{0, 0, 0, 255})
	ctx.Clear()
}

func main() {
	ctx := gg.NewContext(size, size)
	clr(ctx)
	setColor := func(r, g, b uint8) {
		ctx.SetColor(color.RGBA{r, g, b, 255})
	}

	s := func(x float64) int {
		return int(s64 * x)
	}

	grid.Iter().For(func(xy []float64) {
		x, y := xy[0], xy[1]
		sx, sy := s(x), s(y)
		g := dist([]float64{x, y})
		setColor(0, uint8(g), 0)
		ctx.SetPixel(sx, sy)
	})

	s2 := func(x float64) float64 {
		return x * s64
	}

	rand.Seed(time.Now().UnixMicro())
	d := funcs.NewDescender(funcs.Multi{
		Ln: 2,
		M:  dist,
	}, []float64{rand.Float64(), rand.Float64()})
	// Draw each step as it is taken, from red for the first to blue.
	d.Record = func(r funcs.StepRecord) {
		x, y := s2(r.X[0]), s2(r.X[1])
		if r.Step == 0 {
			ctx.MoveTo(x, y)
			return
		}
		c := uint8(255 * math.Min(1, float64(r.Step)/50))
		setColor(255-c, 0, c)
		ctx.LineTo(x, y)
		ctx.Stroke()
		ctx.MoveTo(x, y)
	}
	r := d.Run()
	fmt.Printf("%d steps, %d evaluations, stopped: %s\n", r.Steps, r.Evals, r.Reason)

	ctx.SavePNG("steps.png")

	clr(ctx)
	setColor(255, 0, 0)
	draw(circA, ctx)
	setColor(0, 0, 255)
	draw(circB, ctx)

	xa, ya := circA(r.X[0])
	setColor(255, 255, 0)
	ctx.DrawCircle(xa, ya, 3)
	ctx.Stroke()

	xb, yb := circB(r.X[1])
	setColor(0, 255, 255)
	ctx.DrawCircle(xb, yb, 3)
	ctx.Stroke()

	ctx.SavePNG("result.png")

}

// draw draws the closed curve fn makes over the steps.
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
