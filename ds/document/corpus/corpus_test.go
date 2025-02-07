package corpus_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/document"
	"github.com/adamcolton/luce/ds/document/corpus"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/entity/enttest"
	"github.com/stretchr/testify/assert"
)

func TestRootID(t *testing.T) {
	c := corpus.New()
	words := []string{"this", "is", "a", "test"}
	for i, w := range words {
		assert.Equal(t, corpus.RootID(i), c.WordToID(w))
	}

	for i, w := range words {
		assert.Equal(t, w, c.IDToWord(corpus.RootID(i)))
		assert.Equal(t, corpus.RootID(i), c.WordToID(w))
	}

	assert.Equal(t, "", c.IDToWord(123))

}

func TestVariantID(t *testing.T) {
	c := corpus.New()
	vars := []document.Variant{
		{1, 32},
		{0, 32},
		{1, 13},
		{0, 13},
	}
	for i, v := range vars {
		assert.Equal(t, corpus.VariantID(i), c.VariantToID(v))
	}

	for i, v := range vars {
		id := corpus.VariantID(i)
		assert.Equal(t, v, c.IDToVariant(id))
		assert.Equal(t, id, c.VariantToID(v))
	}

	assert.Equal(t, "", c.IDToWord(123))
}

func TestCorpus(t *testing.T) {
	jabberwocky := `
	'Twas brillig, and the slithy toves
	Did gyre and gimble in the wabe:
	All mimsy were the borogoves,
	And the mome raths outgrabe.
	
	“Beware the Jabberwock, my son!
	The jaws that bite, the claws that catch!
	Beware the Jubjub bird, and shun
	The frumious Bandersnatch!”
	
	He took his vorpal sword in hand;
	Long time the manxome foe he sought—
	So rested he by the Tumtum tree
	And stood awhile in thought.
	
	And, as in uffish thought he stood,
	The Jabberwock, with eyes of flame,
	Came whiffling through the tulgey wood,
	And burbled as it came!
	
	One, two! One, two! And through and through
	The vorpal blade went snicker-snack!
	He left it dead, and with its head
	He went galumphing back.
	
	“And hast thou slain the Jabberwock?
	Come to my arms, my beamish boy!
	O frabjous day! Callooh! Callay!”
	He chortled in his joy.
	
	'Twas brillig, and the slithy toves
	Did gyre and gimble in the wabe:
	All mimsy were the borogoves,
	And the mome raths outgrabe.`
	c := corpus.New()
	d := c.AddDoc(jabberwocky)
	assert.Equal(t, jabberwocky, d.String())
	assert.Equal(t, d, c.GetDoc(d.DocID()))

	assert.Equal(t, d.DocID(), corpus.DocIDer(d).DocID())
}

func TestCorpusSearch(t *testing.T) {
	c := corpus.New()
	lt := slice.LT[string]()
	docs := []*corpus.Document{
		c.AddDoc("The sun was shining on the sea"),
		c.AddDoc("Shining with all it's might"),
		c.AddDoc("And it did the very best it could"),
		c.AddDoc("To make the billows smooth and bright"),
		c.AddDoc("And this was very odd because"),
		c.AddDoc("It was the middle of the night"),
		c.AddDoc("The moon was shining skulkily"),
		c.AddDoc("Because she thought the sun"),
	}
	unfound := c.Find("unfound")
	assert.Nil(t, unfound)

	the := c.Find("the")
	if assert.Equal(t, 6, the.Len()) {
		for _, idx := range []int{0, 2, 3, 5, 6, 7} {
			assert.True(t, the.Contains(docs[idx].DocID()))
		}
	}

	shining := c.Find("shining")
	if assert.Equal(t, 3, shining.Len()) {
		for _, idx := range []int{0, 1, 6} {
			assert.True(t, shining.Contains(docs[idx].DocID()))
		}
	}

	foo := lset.NewMulti(the, shining).Intersection()
	both := c.GetDocs(foo.Slice(nil)).Strings()
	lt.Sort(both)
	expected := []string{
		"The moon was shining skulkily",
		"The sun was shining on the sea",
	}
	assert.Equal(t, expected, both)

	sh := c.Prefix("sh").AllWords().Strings().Slice(nil)
	assert.Equal(t, []string{"she", "shining"}, lt.Sort(sh))

	ffs := c.Containing("ll").AllWords().Strings().Slice(nil)
	assert.Equal(t, []string{"all", "billows"}, lt.Sort(ffs))

}

func TestCorpusEntity(t *testing.T) {
	enttest.Setup()

	jabberwocky := []string{`
	'Twas brillig, and the slithy toves
	Did gyre and gimble in the wabe:
	All mimsy were the borogoves,
	And the mome raths outgrabe.
	`, `
	“Beware the Jabberwock, my son!
	The jaws that bite, the claws that catch!
	Beware the Jubjub bird, and shun
	The frumious Bandersnatch!”
	`, `
	He took his vorpal sword in hand;
	Long time the manxome foe he sought—
	So rested he by the Tumtum tree
	And stood awhile in thought.
	`, `
	And, as in uffish thought he stood,
	The Jabberwock, with eyes of flame,
	Came whiffling through the tulgey wood,
	And burbled as it came!
	`, `
	One, two! One, two! And through and through
	The vorpal blade went snicker-snack!
	He left it dead, and with its head
	He went galumphing back.
	`, `
	“And hast thou slain the Jabberwock?
	Come to my arms, my beamish boy!
	O frabjous day! Callooh! Callay!”
	He chortled in his joy.
	`, `
	'Twas brillig, and the slithy toves
	Did gyre and gimble in the wabe:
	All mimsy were the borogoves,
	And the mome raths outgrabe.`}

	c := corpus.New()
	for _, p := range jabberwocky {
		c.AddDoc(p)
	}

	hisDocs := c.Find("his")
	assert.Equal(t, 2, hisDocs.Len())

	ref, err := c.Save()
	assert.NoError(t, err)

	// This should delete the document and update the corpus record
	id := c.Find("uffish").Slice(nil)[0]
	k := c.GetDoc(id).Key
	assert.True(t, entity.Store.Get(k).Found)
	c.Remove(id)
	assert.Equal(t, 0, c.Find("uffish").Len())

	entity.ClearCache()

	c2 := ref.GetPtr()
	hisDocs = c2.Find("his")
	assert.Equal(t, 2, hisDocs.Len())

	assert.Equal(t, 0, c2.Find("uffish").Len())
}

func TestDocumentUpdate(t *testing.T) {
	c := corpus.New()
	str := "this is a test"
	doc := c.AddDoc(str)
	found := c.Find("is")
	assert.Equal(t, 1, found.Len())
	assert.True(t, found.Contains(doc.ID))
	assert.Equal(t, str, doc.String())

	str = "this was a test"
	// just to prove it's not holding a reference to the string
	assert.NotEqual(t, str, doc.String())
	doc.Update(str)
	assert.Equal(t, str, doc.String())
	found = c.Find("is")
	assert.Equal(t, 0, found.Len())
	assert.False(t, found.Contains(doc.ID))

	found = c.Find("was")
	assert.Equal(t, 1, found.Len())
	assert.True(t, found.Contains(doc.ID))
}

func TestDeleteDocument(t *testing.T) {
	c := corpus.New()
	doc1 := c.AddDoc("this is a test")
	c.AddDoc("this is also a test")
	assert.Equal(t, 2, c.Find("test").Len())
	c.Remove(doc1.ID)
	assert.Equal(t, 1, c.Find("test").Len())
}
