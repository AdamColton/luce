package corpus_test

import (
	"sort"
	"testing"

	"github.com/adamcolton/luce/ds/document/corpus"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/prefix"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

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

	ids := c.Find("gimble")
	assert.True(t, ids.Contains(d.DocID))

	lt := slice.LT[string]()
	bis := c.Prefix("bi").AllWords().Strings().ToSlice(nil)
	assert.Equal(t, []string{"bird", "bite"}, lt.Sort(bis))

	ffs := c.Containing("ff").AllWords().Strings().ToSlice(nil)
	assert.Equal(t, []string{"uffish", "whiffling"}, lt.Sort(ffs))

	expected := []prefix.Suggestion{
		{Word: "te", Terminals: []int{1}},
		{Word: "rd", Terminals: []int{1}},
	}
	assert.Equal(t, expected, c.Prefix("bi").Suggest(10))
}

func TestCorpusSearch(t *testing.T) {
	c := corpus.New()
	c.AddDoc("The sun was shining on the sea")
	c.AddDoc("Shining with all it's might")
	c.AddDoc("And it did the very best it could")
	c.AddDoc("To make the billows smooth and bright")
	c.AddDoc("And this was very odd because")
	c.AddDoc("It was the middle of the night")
	c.AddDoc("The moon was shining skulkily")
	c.AddDoc("Because she thought the sun")

	the := c.Find("the")
	assert.Equal(t, 6, the.Len())

	shining := c.Find("shining")
	assert.Equal(t, 3, shining.Len())

	foo := lset.Multi[corpus.DocID]{the, shining}.Intersection()
	both := c.GetDocs(foo.Slice()).Strings()
	sort.Strings(both)
	expected := []string{
		"The moon was shining skulkily",
		"The sun was shining on the sea",
	}
	assert.Equal(t, expected, both)

	{
		s := c.Prefix("sh").Suggest(10)
		expected := []prefix.Suggestion{
			{
				Word:      "ining",
				Terminals: []int{4},
			},
			{
				Word:      "e",
				Terminals: []int{0},
			},
		}
		assert.Equal(t, expected, s)
	}
}
