package ldate

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/adamcolton/luce/ds/lset"
)

var (
	nowOnce   sync.Once
	nowMu     sync.RWMutex
	now       Date
	listeners = lset.New[*listener]()
)

type listener struct {
	fn       func(Date)
	canceled atomic.Bool
	// callMu keeps the calls to fn in order: the first call, made by
	// DateListener, finishes before the first midnight call starts.
	callMu sync.Mutex
}

// startNow sets now and starts the goroutine that keeps it current, the first
// time it is called.
func startNow() {
	nowOnce.Do(func() {
		n := time.Now()
		setNow(TimeToDate(n))
		go func() {
			waitForTomorrow(n, nil)
			updateNow(time.Now, nil)
		}()
	})
}

// setNow sets now to d and, if that changes the date, calls the listeners.
func setNow(d Date) {
	nowMu.Lock()
	changed := now != d
	now = d
	var ls []*listener
	if changed {
		// Taken in the same lock as setting now, so a listener added after this
		// already got d from DateListener and is not called with it again.
		ls = listeners.Slice(nil)
	}
	nowMu.Unlock()

	for _, l := range ls {
		l.callMu.Lock()
		if !l.canceled.Load() {
			l.fn(d)
		}
		l.callMu.Unlock()
	}
}

// Now returns today's Date in the local time zone. It is updated just after
// midnight.
func Now() Date {
	startNow()
	nowMu.RLock()
	defer nowMu.RUnlock()
	return now
}

// CancelToken stops the calls to a function registered with DateListener.
type CancelToken struct {
	l *listener
}

// Cancel stops the calls. A call already running finishes. Cancel returns
// false if the calls were already canceled, or for a zero CancelToken.
func (c CancelToken) Cancel() bool {
	if c.l == nil || !c.l.canceled.CompareAndSwap(false, true) {
		return false
	}
	nowMu.Lock()
	listeners.Remove(c.l)
	nowMu.Unlock()
	return true
}

// DateListener calls fn with today's Date before it returns, then again just
// after each midnight, with the new Date, until it is canceled. The midnight
// calls run on a goroutine shared by all the listeners, so fn should return
// quickly. It may call Cancel on its own CancelToken.
func DateListener(fn func(Date)) CancelToken {
	startNow()
	l := &listener{fn: fn}
	l.callMu.Lock()
	nowMu.Lock()
	listeners.Add(l)
	d := now
	nowMu.Unlock()
	fn(d)
	l.callMu.Unlock()
	return CancelToken{l: l}
}

// updateNow keeps now current: it sets it from clock, waits for the next day and
// repeats until stop is closed. A nil stop is never closed.
func updateNow(clock func() time.Time, stop <-chan struct{}) {
	for {
		n := clock()
		setNow(TimeToDate(n))
		if !waitForTomorrow(n, stop) {
			return
		}
	}
}

// waitForTomorrow waits until just after the midnight following n. It returns
// false if stop is closed first.
func waitForTomorrow(n time.Time, stop <-chan struct{}) bool {
	tomorrow := time.Date(n.Year(), n.Month(), n.Day()+1, 0, 0, 0, 100, n.Location())
	timer := time.NewTimer(tomorrow.Sub(n))
	select {
	case <-timer.C:
		return true
	case <-stop:
		timer.Stop()
		return false
	}
}
