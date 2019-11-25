package thresher

type HasType interface {
	TypeID() uint64
}

var DefaultSliceThreshold uint64 = 1 << 15
