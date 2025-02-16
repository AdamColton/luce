package entdefertimer

import (
	"time"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/util/dbltimer"
)

type Strategy struct {
	hard, soft time.Duration
	entities   lmap.Wrapper[uint64, *dbltimer.DoubleTimer]
}

func New(hard, soft time.Duration) *Strategy {
	return &Strategy{
		hard:     hard,
		soft:     soft,
		entities: lmap.EmptySafe[uint64, *dbltimer.DoubleTimer](0),
	}
}

func (s *Strategy) DeferSave(er entity.Referer, saveFn func() error) {
	k := er.EntKey()
	h := k.Hash64()
	t, ok := s.entities.Get(h)
	if ok {
		ok = t.Reset()
	}
	if !ok {
		tmr := dbltimer.New(s.hard, s.soft, func() {
			saveFn()
		})
		s.entities.Set(h, tmr)
	}
}

func (s *Strategy) DeferCacheClear(er entity.Referer) {
	// not currently clearing cache
}
