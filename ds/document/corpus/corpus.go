package corpus

// Corpus holds a collection of documents for text indexing. It fulfills
// document.Encoder and document.Decoder.
type Corpus struct {
	cur struct {
		RootID
	}
	id2root   map[RootID]*root
	rootByStr map[string]*root
}

// New creates a Corpus.
func New() *Corpus {
	return &Corpus{
		id2root:   make(map[RootID]*root),
		rootByStr: make(map[string]*root),
	}
}

// WordToID converts a root word to a RootID, fulfilling document.Encoder.
func (c *Corpus) WordToID(rStr string) RootID {
	r := c.rootByStr[rStr]
	if r == nil {
		r = &root{
			RootID: c.cur.RootID,
			str:    rStr,
		}
		c.rootByStr[rStr] = r
		c.id2root[c.cur.RootID] = r
		c.cur.RootID++
	}
	return r.RootID
}

// IDToWord converts a RootID to a root word, fulfilling document.Decoder.
func (c *Corpus) IDToWord(rID RootID) string {
	r := c.id2root[rID]
	if r == nil {
		return ""
	}
	return r.str
}
