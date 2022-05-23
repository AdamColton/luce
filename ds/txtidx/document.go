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
	VarID
	Positions []uint32
}

func (d *Document) String(c *Corpus) string {
	out := make([]byte, d.Ln)
	copy(out, d.start)
	for _, w := range d.Words {
		cw := c.IDs[w.WordID]
		for _, dv := range w.Variants {
			v := c.VariantLookup[dv.VarID]
			vstr := v.apply(cw.Word)

			for _, p := range dv.Positions {
				copy(out[p:], vstr)
			}
		}
	}
	return string(out)
}
