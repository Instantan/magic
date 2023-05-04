package main

import (
	"log"

	"github.com/Instantan/magic"
)

var nameView = magic.View(`
	<span>
		{{name}}
	</span>
`)

var nameComponent = magic.Component(func(s magic.Socket) magic.AppliedView {
	magic.Assign(s, "name", "Child")
	s.HandleEvent(func(ev string, data magic.EventData) {
		log.Println("nameComponent", ev)
	})
	return nameView(s)
})
