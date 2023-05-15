package main

import (
	"github.com/Instantan/magic"
)

var homeView = magic.View(`
	<a href="/overview.html">
		Navigate to overview
	</a>
`)

var home = magic.Component(func(s magic.Socket, e magic.Empty) magic.AppliedView {
	return html(s, HTMLProps{
		title: "Home",
		body:  homeView(s),
	})
})
