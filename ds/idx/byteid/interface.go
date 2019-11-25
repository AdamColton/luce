package byteid

type Index interface {
	Insert(id []byte) (int, bool)
	Get(id []byte) (int, bool)
	Delete(id []byte) (int, bool)
	SliceLen() int
	SetSliceLen(int)
	Next([]byte) ([]byte, int)
}

type IndexFactory func(slicelen int) Index
