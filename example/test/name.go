package main

import (
	"github.com/Instantan/magic"
)

var nameView = magic.View(`
	<span>
		{{name}}
	</span>
`)

var nameComponent = magic.Component(func(s magic.Socket) magic.AppliedView {
	magic.Assign(s, "name", "Child")
	return nameView(s)
})
