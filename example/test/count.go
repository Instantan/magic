package main

import (
	"time"

	"github.com/Instantan/magic"
)

var counterView = magic.View(`
	<h1>
		{{name}} {{ count }}
	</h1>
`)

var counterComponent = magic.Component(func(s magic.Socket) magic.AppliedView {
	c := 0
	magic.Assign(s, "name", nameComponent(s))
	magic.Assign(s, "count", c)

	if s.Live() {
		t := time.NewTicker(time.Second)

		go func() {
			for range t.C {
				c++
				magic.Assign(s, "count", c)
			}
		}()

		s.HandleEvent(func(ev string, data magic.EventData) {
			switch ev {
			case magic.UnmountEvent:
				t.Stop()
			}
		})

	}
	return counterView(s)
})
