package merkle

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
