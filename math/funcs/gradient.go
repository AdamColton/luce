package funcs

// Gradient steps straight down the gradient: plain gradient descent. It is
// simple, but slow when the variables are on very different scales or depend
// on each other.
//
// https://en.wikipedia.org/wiki/Gradient_descent
type Gradient struct {
	// t is the last step length the line search accepted.
	t float64
}

// Reset implements Direction.
func (gr *Gradient) Reset(n int) {
	gr.t = 0
}

// Dir implements Direction.
func (gr *Gradient) Dir(m M, x []float64, f float64, g, d []float64) float64 {
	neg(g, d)
	if gr.t == 0 {
		return 1
	}
	// The line search can only shrink a step, so starting from the last
	// accepted length would let the step shrink but never grow. Starting at
	// four times it lets the step grow again where the function allows it.
	return 4 * gr.t
}

// Update implements Direction.
func (gr *Gradient) Update(s, y []float64, t float64) {
	gr.t = t
}
