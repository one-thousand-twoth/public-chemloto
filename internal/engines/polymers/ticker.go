package polymers

import (
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models/enerr"
)

const (
	Stopped = "Stopped"
	Started = "Started"
)

// Ticker это компонент для state реализующий логику таймера

// No concurrent safe. Use RegisterTo function for binding state
type Ticker struct {
	t      *time.Ticker
	status string

	passed uint
	want   uint
}

type StateAdd interface {
	Add(action string, fun HandlerFunc, secure bool) SimpleState
}

func (t *Ticker) RegisterTo(state StateAdd) {
	state.Add("TimerStop", t.stopTimer(), true)
	state.Add("TimerPlay", t.playTimer(), true)
	state.Add("TimerPause", t.pauseTimer(), true)
}

func (t *Ticker) stopTimer() HandlerFunc {
	return func(e models.Action) (stateInt, error) {
		t.status = Stopped
		t.passed = 0
		return UPDATE_CURRENT, nil
	}
}

func (t *Ticker) playTimer() HandlerFunc {
	return func(e models.Action) (stateInt, error) {
		const op enerr.Op = "TickerComp/playTimer"
		if t.status == Started {
			return NO_TRANSITION, enerr.E(op, "Таймер уже запущен")
		}
		t.status = Started
		return UPDATE_CURRENT, nil
	}
}
func (t *Ticker) pauseTimer() HandlerFunc {
	return func(e models.Action) (stateInt, error) {
		const op enerr.Op = "TickerComp/pauseTimer"
		if t.status == Stopped {
			return NO_TRANSITION, enerr.E(op, "Таймер уже остановлен")
		}
		t.status = Stopped
		return UPDATE_CURRENT, nil
	}
}

// newTicker returns pointer to stopped Ticker.
// Use Ticker.Reset() to start it.
func newTicker() *Ticker {
	t := time.NewTicker(time.Second)
	t.Stop()

	ticker := &Ticker{
		t:      t,
		status: Stopped,
	}
	go func() {
		for range t.C {
			if ticker.status == Started {
				ticker.passed += 1
			}
		}
	}()
	return ticker
}

func (t *Ticker) Status() string {
	return t.status
}
func (t *Ticker) CurrentDuration() time.Duration {
	return time.Duration(t.want) * time.Second
}

func (t *Ticker) Passed() time.Duration {
	return time.Duration(t.passed) * time.Second
}

func (t *Ticker) Remains() time.Duration {
	if t.want <= t.passed {
		return 0
	}
	dur := t.want - t.passed
	return time.Duration(dur) * time.Second
}
func (t *Ticker) Ticked() bool {
	if t.passed >= t.want {
		return true
	}
	return false
}

func (t *Ticker) Stop() {
	t.t.Stop()
	t.status = Stopped
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
	t.passed = 0
	t.want = uint(d.Abs().Seconds())
	t.t.Reset(time.Second)
	t.status = Started
}
