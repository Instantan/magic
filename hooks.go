package magic

type Hook interface {
	Unmount()
}

func Use(s Socket, hook Hook) {
	s.HandleEvent(func(ev string, data EventData) {
		if ev == UnmountEvent {
			hook.Unmount()
		}
	})
}
