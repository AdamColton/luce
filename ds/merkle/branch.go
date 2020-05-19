package merkle

import (
	"io"

	"github.com/adamcolton/luce/lerr"
)

type branch struct {
	children []node
	digest   []byte
	// idx is the idx of the last child
	idx uint32
	// data, count and depth are all cache values. The builder will populate
	// them when building a tree but a tree assembled from validators will not
	// be pre-populated. All three also require a completed tree to be
	// meaningful. A complete tree will only contain branches and dataLeaves,
	// not digestNodes.
	data  []byte
	count uint32
	depth int
	pos   int64
}

func (b *branch) Digest() []byte {
	return b.digest
}

func (b *branch) maxIdx() uint32 {
	return b.idx
}

func (b *branch) Data() []byte {
	if b.data != nil {
		return b.data
	}
	size := b.size()
	if size == -1 {
		return nil
	}
	return b.getData(make([]byte, 0, size))
}

func (b *branch) getData(data []byte) []byte {
	s := len(data)
	for _, c := range b.children {
		if dl, ok := c.(*dataLeaf); ok {
			dlStart := len(data)
			data = append(data, dl.data...)
			// repoint to avoid duplication and free the old memory
			dl.data = data[dlStart:]
		} else {
			data = c.(*branch).getData(data)
		}
	}
	b.data = data[s:]
	return data
}

func (b *branch) size() int {
	sum := 0
	for _, c := range b.children {
		s := c.size()
		if s == -1 {
			return -1
		}
		sum += s
	}
	return sum
}

func (b *branch) Count() uint32 {
	if b.count > 0 {
		return b.count
	}

	var sum uint32
	for _, c := range b.children {
		s := c.Count()
		if s == maxUint32 {
			return s
		}
		sum += s
	}
	b.count = sum
	return sum
}

func (b *branch) Depth() int {
	if b.depth > 0 {
		return b.depth
	}

	depth := 0
	for _, c := range b.children {
		d := c.Depth()
		if d == -1 {
			return d
		}
		if d > depth {
			depth = d
		}
	}
	b.depth = depth + 1
	return b.depth
}

func (b *branch) have(idxs []uint32) []uint32 {
	for _, c := range b.children {
		idxs = c.have(idxs)
	}
	return idxs
}

// Read implements io.Reader
func (b *branch) Read(p []byte) (n int, err error) {
	if b.data == nil {
		b.Data()
	}
	n = len(p)
	l64 := int64(len(b.data))
	if b.pos+int64(n) > l64 {
		n = int(l64 - b.pos)
	}
	if n == 0 {
		err = io.EOF
		return
	}
	end := b.pos + int64(n)
	copy(p, b.data[b.pos:end])
	b.pos = end
	return
}

const (
	// ErrBadWhence is returned if Seek is called and whence is not
	// io.SeekStart, io.SeekEnd or io.SeekCurrent.
	ErrBadWhence = lerr.Str("Bad whence value for Seek")
)

// Seek implements io.Seeker
func (b *branch) Seek(offset int64, whence int) (int64, error) {
	if b.data == nil {
		b.Data()
	}
	switch whence {
	case io.SeekStart:
		b.pos = offset
	case io.SeekEnd:
		b.pos = int64(len(b.data)) + offset
	case io.SeekCurrent:
		b.pos += offset
	default:
		return -1, ErrBadWhence
	}

	l64 := int64(len(b.data))
	if b.pos < 0 {
		b.pos = 0
	} else if b.pos > l64 {
		b.pos = l64
	}

	return b.pos, nil
}
