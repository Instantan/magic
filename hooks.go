package magic

type Hook interface {
	Unmount()
}

// Use is a way to add a hook that listens to an unmount event
func Use(s Socket, hook Hook) {
	s.HandleEvent(func(ev string, _ EventData) {
		if ev == UnmountEvent {
			hook.Unmount()
		}
	})
}

// UseRoutine spawns a new goroutine that gets closed when it receives the unmount event
func UseRoutine(s Socket, fn func(quit <-chan struct{})) {
	q := make(chan struct{})
	s.HandleEvent(func(ev string, _ EventData) {
		if ev == UnmountEvent {
			q <- struct{}{}
		}
	})
	go fn(q)
}

// UseLiveRoutine spawns a new goroutine when the connection is live
func UseLiveRoutine(s Socket, fn func(quit <-chan struct{})) {
	if s.Live() {
		UseRoutine(s, fn)
	}
}
