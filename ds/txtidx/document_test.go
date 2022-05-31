package txtidx

import (
	"io/ioutil"
	"sort"
	"testing"

	"github.com/adamcolton/luce/ds/huffman"
	"github.com/adamcolton/luce/serial/rye"
	"github.com/stretchr/testify/assert"
)

func TestDocumentString(t *testing.T) {
	expected := "--- This is a Test "

	c := NewCorpus()

	d := c.AddDoc(expected)

	assert.Equal(t, expected, c.DocString(d))

	// 4 words but only 2 variants
	// "is " and "a " have the same variant
	// "This " and "Test " have the same variant
	assert.Len(t, c.variantsByStr, 2)

	assert.True(t, c.find("test").Has(d))
}

func TestMultiDoc(t *testing.T) {
	c := NewCorpus()
	c.AddDoc("The sun was shining on the sea")
	c.AddDoc("Shining with all it's might")
	c.AddDoc("And it did the very best it could")
	c.AddDoc("To make the billows smooth and bright")
	c.AddDoc("And this was very odd because")
	c.AddDoc("It was the middle of the night")
	c.AddDoc("The moon was shining skulkily")
	c.AddDoc("Because she thought the sun")

	the := c.find("the")
	assert.Equal(t, 6, the.Len())

	shining := c.find("shining")
	assert.Equal(t, 3, shining.Len())

	both := c.GetDocs(c.find("the", "shining"))
	sort.Strings(both)
	expected := []string{
		"The moon was shining skulkily",
		"The sun was shining on the sea",
	}
	assert.Equal(t, expected, both)
	_, both = c.Search("the shining")
	sort.Strings(both)
	assert.Equal(t, expected, both)

	ill := c.GetDocs(c.find("ill"))
	expected = []string{
		"To make the billows smooth and bright",
	}
	assert.Equal(t, expected, ill)

	_, exact := c.Search(`"The sun"`)
	expected = []string{
		"The sun was shining on the sea",
	}
	assert.Equal(t, expected, exact)

	ll := c.GetDocs(c.find("ll"))
	expected = []string{
		"Shining with all it's might",
		"To make the billows smooth and bright",
	}
	assert.Equal(t, expected, ll)
}

func TestMarkov(t *testing.T) {
	m := newMarkov()
	w := m.upsert("test")
	id := DocID(123)
	w.Documents.Add(id)
	w2 := m.upsert("test")

	assert.True(t, w2.Documents.Has(id))

	w3, _ := m.find("test")
	assert.True(t, w3.Documents.Has(id))
}

func TestDelete(t *testing.T) {
	c := NewCorpus()
	d0 := c.AddDoc("this is document 0 keyphrase")
	assert.Equal(t, DocID(0), d0.ID())

	d1 := c.AddDoc("this is document 1")
	assert.Equal(t, DocID(1), d1.ID())

	w, _ := c.roots.find("keyphrase")
	assert.True(t, w.Documents.Has(d0.ID()))
	wid := w.wordIDX

	c.Delete(d0.ID())
	assert.Nil(t, c.docs[0])
	assert.Nil(t, c.words[wid])

	w, _ = c.roots.find("keyphrase")
	assert.Nil(t, w)

	d2 := c.AddDoc("this is document 2, reallocated DocIDX 0 keyphrase")
	assert.Equal(t, DocID(0), d2.ID())
	assert.NotNil(t, c.words[wid])

	w, _ = c.roots.find("keyphrase")
	assert.True(t, w.Documents.Has(d2.ID()))

	str := "this is document 2.1 - it has been updated"
	c.Update(d2, str)
	assert.Equal(t, c.DocString(d2), str)
	w, _ = c.roots.find("keyphrase")
	assert.Nil(t, w)
}

func TestCompression(t *testing.T) {
	f, err := ioutil.ReadFile("/home/adam/Projects/homestead/workshop/workshop.wiki")
	assert.NoError(t, err)

	c := NewCorpus()
	hd := newDoc(string(f), c)
	str := hd.toString(c)
	assert.Equal(t, string(f), str)

	size := hd.vEnc.Ln + hd.wEnc.Ln

	wl := huffman.NewLookup(hd.wt)
	vl := huffman.NewLookup(hd.vt)

	for _, wIdx := range wl.All() {
		size += int(rye.SizeCompactUint64(uint64(wIdx))) * 8
		size += wl.Get(wIdx).Ln
	}
	for _, vIdx := range vl.All() {
		size += int(rye.SizeCompactUint64(uint64(vIdx))) * 8
		size += vl.Get(vIdx).Ln
	}

	size = (size / 8) + len(hd.wSingles) + len(hd.vSingles)

	assert.Equal(t, len(f), size)

}
