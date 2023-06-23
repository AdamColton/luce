package corpus_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/document"
	"github.com/adamcolton/luce/ds/document/corpus"
	"github.com/adamcolton/luce/ds/slice"
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
	str := `
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
	d := c.AddDoc(str)
	assert.Equal(t, str, d.String())

	assert.Equal(t, d.DocID, corpus.DocIDer(d).ID())
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
			assert.True(t, the.Contains(docs[idx].ID()))
		}
	}

	shining := c.Find("shining")
	if assert.Equal(t, 3, shining.Len()) {
		for _, idx := range []int{0, 1, 6} {
			assert.True(t, shining.Contains(docs[idx].ID()))
		}
	}

	sh := c.Prefix("sh").AllWords().Strings().ToSlice(nil)
	assert.Equal(t, []string{"she", "shining"}, lt.Sort(sh))

	ffs := c.Containing("ll").AllWords().Strings().ToSlice(nil)
	assert.Equal(t, []string{"all", "billows"}, lt.Sort(ffs))
}
