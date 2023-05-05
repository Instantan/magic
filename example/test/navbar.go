package main

import (
	"time"

	"github.com/Instantan/magic"
)

var navbarView = magic.View(`
	<nav magic:click="test">
		{{ content }}
	</nav>
`)

var navbarComponent = magic.Component(func(s magic.Socket) magic.AppliedView {
	magic.Assign(s, "content", time.Now().Local().Format(time.RFC1123))
	if s.Live() {
		t := time.NewTicker(time.Second)
		go func() {
			for range t.C {
				magic.Assign(s, "content", time.Now().Local().Format(time.RFC1123))
			}
		}()
		s.HandleEvent(func(ev string, data magic.EventData) {
			switch ev {
			case magic.ClickEvent:
				print(string(data))
			case magic.UnmountEvent:
				t.Stop()
			}
		})
	}
	return navbarView(s)
})