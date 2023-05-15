package main

import "github.com/Instantan/magic"

var overviewView = magic.View(`
	<a href="/index.html">
		Navigate to home
	</a>
`)

var overview = magic.Component(func(s magic.Socket, e magic.Empty) magic.AppliedView {
	return html(s, HTMLProps{
		title: "Home",
		body:  overviewView(s),
	})
})
