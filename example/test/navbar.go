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
		_ = t
		go func() {
			for {
				select {
				case <-t.C:
					magic.Assign(s, "content", time.Now().Local().Format(time.RFC1123))
				}
			}
		}()
	}
	return navbarView(s)
})
