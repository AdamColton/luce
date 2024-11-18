package ljson

// MarshalString fulfills Marshaler for a string
func MarshalString[Ctx any](str string, ctx *MarshalContext[Ctx]) (WriteNode, error) {
	return func(ctx *WriteContext) {
		ctx.Cache = EncodeString(ctx.Cache[:0], str, ctx.EscapeHtml)
		ctx.FlushCache()
	}, nil
}
