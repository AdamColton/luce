package document

import (
	"strings"
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

type wordID uint32
type varID uint32

type mockEncDec struct {
	word2id map[string]wordID
	id2Word []string
	var2id  map[string]varID
	id2Var  []Variant
}

func (med *mockEncDec) WordToID(word string) wordID {
	id, found := med.word2id[word]
	if !found {
		id = wordID(len(med.id2Word))
		med.word2id[word] = id
		med.id2Word = append(med.id2Word, word)
	}
	return id
}
func (med *mockEncDec) VariantToID(v Variant) varID {
	id, found := med.var2id[string(v)]
	if !found {
		id = varID(len(med.id2Var))
		med.var2id[string(v)] = id
		med.id2Var = append(med.id2Var, v)
	}
	return id
}
func (med *mockEncDec) IDToWord(id wordID) string {
	return med.id2Word[id]
}
func (med *mockEncDec) IDToVariant(id varID) Variant {
	return med.id2Var[id]

}

func TestDocument2(t *testing.T) {
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
	med := &mockEncDec{
		word2id: map[string]wordID{},
		var2id:  map[string]varID{},
	}

	maxUint32 := ^uint32(0)

	enc := &DocumentEncoder[wordID, varID]{
		Encoder:         med,
		Splitter:        Parse,
		RootVariant:     RootVariant,
		WordSingleToken: wordID(maxUint32),
		VarSingleToken:  varID(maxUint32),
	}

	doc := enc.Build(str)

	dec := &DocumentDecoder[wordID, varID]{
		Decoder:         med,
		WordSingleToken: enc.WordSingleToken,
		VarSingleToken:  enc.VarSingleToken,
	}

	got := dec.Decode(doc)

	assert.Equal(t, str, got)
	lt := slice.LT[wordID]()

	{
		// words should return one instance of every index value because this is
		// not a shared corpus.
		expected := make([]wordID, len(med.word2id))
		for i := range expected {
			expected[i] = wordID(i)
		}
		words := lt.Sort(doc.WordIDs())
		assert.Equal(t, expected, words)
	}

	{
		str = strings.NewReplacer("brillig", "gloorp", "borogoves", "klatuu").Replace(str)
		cs := enc.Update(doc, str)
		lt.Sort(cs.Added)
		lt.Sort(cs.Removed)
		expected := &ChangeSet[wordID]{
			Removed: []wordID{med.WordToID("brillig"), med.WordToID("borogoves")},
			Added:   []wordID{med.WordToID("gloorp"), med.WordToID("klatuu")},
		}
		assert.Equal(t, expected, cs)
	}
}
