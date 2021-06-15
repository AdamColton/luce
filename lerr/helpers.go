package lerr

// Panic if err is not nil. If err is in the exception list, it will return
// true, but will not panic.
func Panic(err error, except ...error) bool {
	for _, ex := range except {
		if err == ex {
			return true
		}
	}
	if err != nil {
		panic(err)
	}
	return false
}
