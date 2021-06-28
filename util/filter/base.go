package filter

// Int provides tools to filter ints and compose filters
type Int func(int) bool

// Returns all values that return true when passed to Int.
func (i Int) Slice(vals []int) []int {
	var out []int
	for _, val := range vals {
		if i(val) {
			out = append(out, val)
		}
	}
	return out
}

// Chan runs a go routine listening on ch and any int that passes the
// Int is passed to the channel that is returned.
func (i Int) Chan(ch <-chan int, buf int) <-chan int {
	out := make(chan int, buf)
	go func() {
		for in := range ch {
			if i(in) {
				out <- in
			}
		}
		close(out)
	}()
	return out
}

// Or builds a new Int that will return true if either underlying
// Int is true.
func (i Int) Or(i2 Int) Int {
	return func(val int) bool {
		return i(val) || i2(val)
	}
}

// And builds a new Int that will return true if both underlying
// Ints are true.
func (i Int) And(i2 Int) Int {
	return func(val int) bool {
		return i(val) && i2(val)
	}
}

// Not builds a new Int that will return true if the underlying
// Int is false.
func (i Int) Not() Int {
	return func(val int) bool {
		return !i(val)
	}
}

// String provides tools to filter strings and compose filters
type String func(string) bool

// Returns all values that return true when passed to String.
func (s String) Slice(vals []string) []string {
	var out []string
	for _, val := range vals {
		if s(val) {
			out = append(out, val)
		}
	}
	return out
}

// Chan runs a go routine listening on ch and any string that passes the
// String is passed to the channel that is returned.
func (s String) Chan(ch <-chan string, buf int) <-chan string {
	out := make(chan string, buf)
	go func() {
		for in := range ch {
			if s(in) {
				out <- in
			}
		}
		close(out)
	}()
	return out
}

// Or builds a new String that will return true if either underlying
// String is true.
func (s String) Or(s2 String) String {
	return func(val string) bool {
		return s(val) || s2(val)
	}
}

// And builds a new String that will return true if both underlying
// Strings are true.
func (s String) And(s2 String) String {
	return func(val string) bool {
		return s(val) && s2(val)
	}
}

// Not builds a new String that will return true if the underlying
// String is false.
func (s String) Not() String {
	return func(val string) bool {
		return !s(val)
	}
}

// Float provides tools to filter float64s and compose filters
type Float func(float64) bool

// Returns all values that return true when passed to Float.
func (f Float) Slice(vals []float64) []float64 {
	var out []float64
	for _, val := range vals {
		if f(val) {
			out = append(out, val)
		}
	}
	return out
}

// Chan runs a go routine listening on ch and any float64 that passes the
// Float is passed to the channel that is returned.
func (f Float) Chan(ch <-chan float64, buf int) <-chan float64 {
	out := make(chan float64, buf)
	go func() {
		for in := range ch {
			if f(in) {
				out <- in
			}
		}
		close(out)
	}()
	return out
}

// Or builds a new Float that will return true if either underlying
// Float is true.
func (f Float) Or(f2 Float) Float {
	return func(val float64) bool {
		return f(val) || f2(val)
	}
}

// And builds a new Float that will return true if both underlying
// Floats are true.
func (f Float) And(f2 Float) Float {
	return func(val float64) bool {
		return f(val) && f2(val)
	}
}

// Not builds a new Float that will return true if the underlying
// Float is false.
func (f Float) Not() Float {
	return func(val float64) bool {
		return !f(val)
	}
}
