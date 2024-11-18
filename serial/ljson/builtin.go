package ljson

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

// MarshalString fulfills Marshaler for a string
func MarshalString[Ctx any](str string, ctx *MarshalContext[Ctx]) (WriteNode, error) {
	return func(ctx *WriteContext) {
		ctx.Cache = EncodeString(ctx.Cache[:0], str, ctx.EscapeHtml)
		ctx.FlushCache()
	}, nil
}

const (
	tstr = "true"
	fstr = "false"
)

// MarshalBool fulfills Marshaler for a bool.
func MarshalBool[Ctx any](b bool, ctx *MarshalContext[Ctx]) (WriteNode, error) {
	str := tstr
	if !b {
		str = fstr
	}
	return func(ctx *WriteContext) {
		ctx.WriteString(str)
	}, nil
}

// MarshalBool fulfills Marshaler for any signed int.
func MarshalInt[I constraints.Signed, Ctx any](i I, ctx *MarshalContext[Ctx]) (WriteNode, error) {
	return func(ctx *WriteContext) {
		ctx.WriteString(strconv.FormatInt(int64(i), 10))
	}, nil
}

// MarshalBool fulfills Marshaler for any unsigned int.
func MarshalUint[U constraints.Unsigned, Ctx any](u U, ctx *MarshalContext[Ctx]) (WriteNode, error) {
	return func(ctx *WriteContext) {
		ctx.WriteString(strconv.FormatUint(uint64(u), 10))
	}, nil
}

// MarshalBool fulfills Marshaler for any float.
func MarshalFloat[F constraints.Float, Ctx any](f F, ctx *MarshalContext[Ctx]) (WriteNode, error) {
	return func(ctx *WriteContext) {
		ctx.WriteString(strconv.FormatFloat(float64(f), 'g', -1, 64))
	}, nil
}
