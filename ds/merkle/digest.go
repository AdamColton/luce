package merkle

type digestNode []byte

func (d digestNode) Digest() []byte            { return d }
func (digestNode) private()                    {}
func (digestNode) size() int                   { return -1 }
func (digestNode) Count() int                  { return -1 }
func (digestNode) Depth() int                  { return -1 }
func (digestNode) maxIdx() uint32              { return ^uint32(0) }
func (digestNode) have(idxs []uint32) []uint32 { return idxs }
