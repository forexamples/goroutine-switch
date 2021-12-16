package ticks

import (
	"log"
	"time"
)

type TickerI interface {
	Tick(cb func())
	Stop()
}

type Ticker struct {
	d        time.Duration
	t        *time.Ticker
	doneChan chan struct{}
	cb       func()
	running  bool
}

func NewTicker(d time.Duration) TickerI {
	ret := &Ticker{
		d: d,
	}
	return ret
}

func (t *Ticker) Tick(cb func()) {
	t.t = time.NewTicker(t.d)
	t.doneChan = make(chan struct{})
	t.cb = cb
	if t.running {
		log.Printf("ticker is running")
		return
	}

	t.running = true

	go func() {
		defer t.t.Stop()
		for {
			select {
			case <-t.doneChan:
				return
			case <-t.t.C:
				cb()
			}
		}
	}()
}

func (t *Ticker) Stop() {
	if !t.running {
		return
	}

	t.running = false
	close(t.doneChan)
}
