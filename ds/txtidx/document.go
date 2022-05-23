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

type DocSet struct {
	docs map[DocID]sig
}

func newDocSet() *DocSet {
	return &DocSet{
		docs: map[DocID]sig{},
	}
}

func (ds *DocSet) add(di DocID) {
	ds.docs[di] = sig{}
}

func (ds *DocSet) Has(di DocID) bool {
	_, found := ds.docs[di]
	return found
}

func (ds *DocSet) Len() int {
	return len(ds.docs)
}

func (ds *DocSet) Intersect(with *DocSet) *DocSet {
	out := map[DocID]sig{}
	iter, cmpr := ds.docs, with.docs
	if len(iter) > len(cmpr) {
		iter, cmpr = cmpr, iter
	}
	for di := range iter {
		_, found := cmpr[di]
		if found {
			out[di] = sig{}
		}
	}
	return &DocSet{
		docs: out,
	}
}

func (ds *DocSet) Slice(c *Corpus) []string {
	out := make([]string, 0, len(ds.docs))
	for di := range ds.docs {
		out = append(out, c.Docs[di].String(c))
	}
	return out
}
