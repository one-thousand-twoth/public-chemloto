package polymers

import (
	"sync"
	"time"
)

type Ticker struct {
	t               *time.Ticker
	currentDuration time.Duration
	startTime       time.Time
	mu              sync.Mutex
}

// newTicker returns pointer to stoped Ticker.
// Use Ticker.Reset() to start it.
func newTicker() *Ticker {
	t := time.NewTicker(time.Hour)
	t.Stop()

	ticker := &Ticker{
		t: t,
	}
	return ticker
}

func (t *Ticker) CurrentDuration() time.Duration {
	return t.currentDuration
}

func (t *Ticker) StartTime() time.Time {
	return t.startTime
}

func (t *Ticker) Passed() time.Duration {
	return time.Since(t.startTime)
}

func (t *Ticker) Remains() time.Duration {
	dur := t.currentDuration - time.Since(t.startTime)
	if dur < 0 {
		return 0
	}
	return dur
}
func (t *Ticker) Ticked() bool {
	select {
	case <-t.t.C:
		return true
	default:
		return false
	}
}

func (t *Ticker) Stop() {
	t.t.Stop()
	select {
	case <-t.t.C:
		// drain ticker channel
	default:
		// fmt.Println("drain nothing")
	}

}
func (t *Ticker) ResetWithFunction(d time.Duration) {
	t.Reset(d)
}
func (t *Ticker) Reset(d time.Duration) {
	t.Stop()
	t.currentDuration = d
	t.startTime = time.Now()
	t.t.Reset(d)
}

func (t *Ticker) Lock(d time.Duration) {
	t.mu.Lock()
}
func (t *Ticker) UnLock(d time.Duration) {
	t.mu.Unlock()
}
