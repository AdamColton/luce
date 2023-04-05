package txtidx

import (
	"github.com/adamcolton/luce/ds/huffman"
	"github.com/adamcolton/luce/serial/rye"
	"github.com/adamcolton/luce/util/lstr"
)

type Document interface {
	Len() int
	ID() DocID
}

type document struct {
	id                 DocID
	ln                 int
	start              []byte
	wt                 huffman.Tree[wordIDX]
	vt                 huffman.Tree[varIDX]
	wSingles, vSingles []byte
	wEnc, vEnc         *rye.Bits
}

func (d *document) ID() DocID {
	return d.id
}

func (d *document) Len() int {
	return d.ln
}

type docUnit[T unit] struct {
	counts     map[T]int
	singles    map[T]sig
	idxs       []T
	singlesDat []byte
}

func newHuffDocUnit[T unit](ln int) *docUnit[T] {
	return &docUnit[T]{
		counts:  make(map[T]int),
		singles: make(map[T]sig),
		idxs:    make([]T, 0, ln),
	}
}

type unit interface {
	~uint32
}

func (hdu *docUnit[T]) add(u T) {
	hdu.counts[u]++
	hdu.idxs = append(hdu.idxs, u)
}

func (hdu *docUnit[T]) encode() (huffman.Tree[T], *rye.Bits, []byte) {
	var singleToken T = T(MaxUint32)
	for u, c := range hdu.counts {
		if c == 1 {
			hdu.singles[u] = sig{}
		}
	}
	for w := range hdu.singles {
		delete(hdu.counts, w)
	}
	hdu.counts[singleToken] = len(hdu.singles)

	s := &rye.Serializer{}
	for i, u := range hdu.idxs {
		_, isSingle := hdu.singles[u]
		if isSingle {
			hdu.idxs[i] = singleToken
			s.CheckFree(5)
			s.CompactUint64(uint64(u))
		}
	}
	hdu.singlesDat = s.Data[:s.Idx]

	if len(hdu.counts) < 2 {
		return nil, nil, hdu.singlesDat
	}

	t := huffman.MapNew(hdu.counts)
	return t, huffman.Encode(hdu.idxs, huffman.NewLookup(t)), hdu.singlesDat
}

func newDoc(str string, c *Corpus) *document {
	start, words := parse(lstr.NewScanner(str))

	wUnit := newHuffDocUnit[wordIDX](len(words))
	vUnit := newHuffDocUnit[varIDX](len(words))

	for _, w := range words {
		wIdx, vIdx := c.upsert(w)
		wUnit.add(wIdx)
		vUnit.add(vIdx)
	}

	hd := &document{
		ln:    len(str),
		start: start,
	}
	c.allocDocIDX(hd)
	for wIdx := range wUnit.counts {
		c.words[wIdx].Documents.Add(hd)
	}

	hd.wt, hd.wEnc, hd.wSingles = wUnit.encode()
	hd.vt, hd.vEnc, hd.vSingles = vUnit.encode()

	return hd
}

func (hd *document) decode() (wIdxs []wordIDX, vIdxs []varIDX) {
	wm := wordIDX(MaxUint32)
	vm := varIDX(MaxUint32)

	dw := rye.NewDeserializer(hd.wSingles)
	if hd.wt == nil {
		ln := len(dw.Data)
		for dw.Idx < ln {
			wIdxs = append(wIdxs, wordIDX(dw.CompactUint64()))
		}
	} else {
		wIdxs = hd.wt.ReadAll(hd.wEnc.Reset())
		for i, w := range wIdxs {
			if w == wm {
				wIdxs[i] = wordIDX(dw.CompactUint64())
			}
		}
	}

	dv := rye.NewDeserializer(hd.vSingles)
	if hd.vt == nil {
		ln := len(dv.Data)
		for dv.Idx < ln {
			vIdxs = append(vIdxs, varIDX(dv.CompactUint64()))
		}
	} else {
		vIdxs = hd.vt.ReadAll(hd.vEnc.Reset())
		for i, v := range vIdxs {
			if v == vm {
				vIdxs[i] = varIDX(dv.CompactUint64())
			}
		}
	}

	return wIdxs, vIdxs
}

func (hd *document) toString(c *Corpus) string {
	wIdxs, vIdxs := hd.decode()
	out := make([]byte, 0, hd.ln)
	out = append(out, hd.start...)
	for i, v := range vIdxs {
		s := c.variants[v].apply(c.words[wIdxs[i]].str)
		out = append(out, s...)
	}
	return string(out)
}

func (hd *document) words() (wIdxs []wordIDX) {
	if hd.wt != nil {
		wIdxs = huffman.NewLookup(hd.wt).All()
	}

	d := rye.NewDeserializer(hd.wSingles)
	ln := len(d.Data)
	for d.Idx < ln {
		wIdxs = append(wIdxs, wordIDX(d.CompactUint64()))
	}
	return wIdxs
}

func (hd *document) update(c *Corpus, str string) {
	wIdxs := hd.words()
	wordsBefore := make(map[wordIDX]sig, len(wIdxs))
	for _, wIdx := range wIdxs {
		wordsBefore[wIdx] = sig{}
	}
	start, words := parse(lstr.NewScanner(str))
	hd.start = start

	wUnit := newHuffDocUnit[wordIDX](len(words))
	vUnit := newHuffDocUnit[varIDX](len(words))

	for _, w := range words {
		wIdx, vIdx := c.upsert(w)
		wUnit.add(wIdx)
		vUnit.add(vIdx)
	}

	hd.wt, hd.wEnc, hd.wSingles = wUnit.encode()
	hd.vt, hd.vEnc, hd.vSingles = vUnit.encode()

	// this could be done directly instead of building the intermediary slice
	wIdxs = hd.words()
	for _, wIdx := range wIdxs {
		_, found := wordsBefore[wIdx]
		if found {
			delete(wordsBefore, wIdx)
		} else {
			c.words[wIdx].Documents.Add(hd)
		}
	}
	for wIdx := range wordsBefore {
		c.deleteDocWord(hd, c.words[wIdx])
	}
}
