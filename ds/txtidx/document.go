package txtidx

type Document struct {
	Ln uint32
	DocID
	Words []DocWord
	start []byte
}

type DocID uint32

type DocWord struct {
	WordID
	Variants []DocVar
}

type DocVar struct {
	VIDX
	Positions []uint32
}

func (d *Document) String(c *Corpus) string {
	out := make([]byte, d.Ln)
	copy(out, d.start)
	for _, w := range d.Words {
		cw := c.IDs[w.WordID]
		for _, v := range w.Variants {
			b := cw.Variants[v.VIDX]
			for _, p := range v.Positions {
				copy(out[p:], b)
			}
		}
	}
	return string(out)
}
