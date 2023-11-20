package funcs

// Direction picks the direction of each step of a Descender. A Direction can
// keep state between steps, so each Descender needs its own.
type Direction interface {
	// Reset starts again for n variables, forgetting any history.
	Reset(n int)
	// Dir writes the direction to step from x into d, where m has the value f
	// and the gradient g, and returns the step length for the line search to
	// try first. The direction should lead down (d·g < 0); if it doesn't, the
	// Descender uses -g and calls Reset. Dir may call m, and change x while
	// it does, as long as x is restored before Dir returns. Those calls count
	// as evaluations.
	Dir(m M, x []float64, f float64, g, d []float64) (t float64)
	// Update is told about each step the line search accepted: s is the
	// change in x, y the change in the gradient and t the step length used.
	Update(s, y []float64, t float64)
}
