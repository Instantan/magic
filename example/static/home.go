package main

import (
	"time"

	"github.com/Instantan/magic"
)

var homeView = magic.View(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<meta name="theme-color" content="#35B6D2" />
		<title>Script</title>
	</head>
	<style>
		html {
			background-color: black;
			color: rgb(200, 200, 200);
		}
	</style>
	<body magic:static="a">
		<p>{{time}}</p>
	</body>
	</html>
`)

var home = magic.Component(func(s magic.Socket, e magic.Empty) magic.AppliedView {
	magic.UseLiveRoutine(s, func(quit <-chan struct{}) {
		t := time.NewTicker(time.Millisecond * 1000)
		for {
			select {
			case c := <-t.C:
				magic.Assign(s, "time", c.String())
			case <-quit:
			}
		}
	})
	return homeView(s)
})
