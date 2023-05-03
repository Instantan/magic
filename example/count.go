package main

import (
	"time"

	"github.com/Instantan/magic"
)

var counterView = magic.View(`
	<h1>
		{{ count }}
	</g1>
`)

var counterComponent = magic.Component(func(s magic.Socket) magic.AppliedView {
	c := 0
	magic.Assign(s, "count", c)
	if s.Live() {
		t := time.NewTicker(time.Second)
		go func() {
			for {
				select {
				case <-t.C:
					c++
					magic.Assign(s, "count", c)
				case <-s.Done():
					t.Stop()
					return
				}
			}
		}()
	}
	return counterView(s)
})
