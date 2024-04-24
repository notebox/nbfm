package combine

import (
	"sync"
	"time"
)

type Debounce struct {
	mu    sync.Mutex
	after time.Duration
	timer *time.Timer
}

func NewDebounce(after time.Duration) *Debounce {
	return &Debounce{after: after}
}

func (d *Debounce) Add(f func()) (count int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		count -= 1
	}
	d.timer = time.AfterFunc(d.after, func() {
		d.mu.Lock()
		d.timer = nil
		d.mu.Unlock()
		f()
	})
	count += 1
	return
}

func (d *Debounce) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Reset(0)
	}
}

func (d *Debounce) stop() (count int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
		count -= 1
	}
	return
}

type Throttle struct {
	mu       sync.Mutex
	duration time.Duration
	ignore   bool
}

func NewThrottle(duration time.Duration) *Throttle {
	return &Throttle{duration: duration}
}

func (t *Throttle) Add() (count int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.ignore {
		return
	}
	t.ignore = true

	time.AfterFunc(t.duration, func() {
		t.mu.Lock()
		defer t.mu.Unlock()
		t.ignore = false
	})

	return 1
}

type DebouncingThrottle struct {
	Throttle
	Debounce
}

func NewDebouncingThrottle(duration time.Duration) *DebouncingThrottle {
	return &DebouncingThrottle{
		Throttle: *NewThrottle(duration),
		Debounce: *NewDebounce(duration),
	}
}

func (td *DebouncingThrottle) Add(fn func()) (count int) {
	count = td.Throttle.Add()
	if count == 1 {
		count += td.Debounce.stop()
		fn()
		return
	}
	return td.Debounce.Add(fn)
}

func (td *DebouncingThrottle) Clear() {
	td.Debounce.Clear()
}

func (td *DebouncingThrottle) Stop() {
	td.Debounce.stop()
}
