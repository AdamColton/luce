package gothicgo

// MemoryContext provides a Context that will not write to disk. Primarily
// intended for testing.
func MemoryContext() Context {
	ctx := CtxFactory{}.New()

	return ctx
}
