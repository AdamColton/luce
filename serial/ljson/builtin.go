package ljson

// MarshalString converts a string to a WriteNode
func MarshalString(str string) WriteNode {
	return func(ctx *WriteContext) {
		ctx.Cache = EncodeString(ctx.Cache[:0], str, ctx.EscapeHtml)
		ctx.FlushCache()
	}
}
