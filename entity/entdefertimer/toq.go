package entdefertimer

import (
	"sync"
	"time"

	"github.com/adamcolton/luce/ds/toq"
	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/lerr"
)

type TOQ struct {
	saveQ    *toq.TimeoutQueue
	clearQ   *toq.TimeoutQueue
	toSave   map[uint64]toq.Token
	toClear  map[uint64]toq.Token
	saveMux  sync.Mutex
	clearMux sync.Mutex
	done     bool
	C        int
}

const ErrClearTooShort = lerr.Str("clear must be at least 50% longer than save or zero")

func NewToq(clear, save time.Duration, capacity int) *TOQ {
	if clear > 0 && 2*clear < 3*save {
		panic(ErrClearTooShort)
	}
	q := &TOQ{
		saveQ:  toq.New(save, capacity),
		toSave: make(map[uint64]toq.Token),
		done:   true,
	}

	if clear > 0 {
		q.toClear = make(map[uint64]toq.Token)
		q.clearQ = toq.New(clear, capacity)
	}

	return q
}

func (t *TOQ) DeferSave(er entity.Referer, saveFn func() error) {
	t.done = false
	k := er.EntKey()
	h := k.Hash64()
	t.saveMux.Lock()
	tok, ok := t.toSave[h]
	if ok {
		ok = tok.Reset()
	}
	if !ok {
		action := func() {
			saveFn()
			t.saveMux.Lock()
			delete(t.toSave, h)
			t.done = len(t.toSave) == 0
			t.saveMux.Unlock()
		}
		t.toSave[h] = t.saveQ.Add(action)
	}
	t.saveMux.Unlock()
}

// Done returns true when all saves are complete
func (t *TOQ) Done() bool {
	return t.done
}

// Flush forces all saves to happen immediatly
func (t *TOQ) Flush() {
	// todo: operate in batches with brief unlocks to allow new saves to
	// come in
	panic("not implemented")
}

func (t *TOQ) DeferCacheClear(er entity.Referer) {
	if t.clearQ != nil {
		k := er.EntKey()
		h := k.Hash64()
		t.clearMux.Lock()
		tok, ok := t.toClear[h]
		if ok {
			ok = tok.Reset()
		}
		if !ok {
			action := func() {
				er.Clear(true)
				t.clearMux.Lock()
				delete(t.toClear, h)
				t.clearMux.Unlock()
			}
			t.toClear[h] = t.clearQ.Add(action)
		}
		t.clearMux.Unlock()
	}
}
