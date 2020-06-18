package lerr

// Panic if err is not nil
func Panic(err error) {
	if err != nil {
		panic(err)
	}
}
