package entdefertimer

import (
	"sync"
	"time"

	"github.com/adamcolton/luce/ds/toq"
	"github.com/adamcolton/luce/entity"
)

type TOQ struct {
	q        *toq.TimeoutQueue
	entities map[uint64]toq.Token
	mux      sync.Mutex
	done     bool
	C        int
}

func NewToq(d time.Duration) *TOQ {
	return &TOQ{
		q:        toq.New(d, 10000),
		entities: make(map[uint64]toq.Token),
		done:     true,
	}
}

func (t *TOQ) DeferSave(er entity.Referer, saveFn func() error) {
	t.done = false
	k := er.EntKey()
	h := k.Hash64()
	t.mux.Lock()
	tok, ok := t.entities[h]
	if ok {
		ok = tok.Reset()
	}
	if !ok {
		action := func() {
			saveFn()
			t.mux.Lock()
			delete(t.entities, h)
			t.done = len(t.entities) == 0
			t.mux.Unlock()
		}
		t.entities[h] = t.q.Add(action)
	}
	t.mux.Unlock()
}

func (t *TOQ) Done() bool {
	return t.done
}

func (t *TOQ) DeferCacheClear(er entity.Referer) {
	// not currently clearing cache
}
