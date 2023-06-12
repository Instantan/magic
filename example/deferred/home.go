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
		<img id="test" src="https://picsum.photos/200/300" />
		{{time}}
	</body>
	</html>
`)

var testView = magic.View(`
	<p>{{child}}</p>
`)

var test = magic.Component(func(s magic.Socket, e magic.Empty) magic.AppliedView {
	return testView(s)
})

var home = magic.Component(func(s magic.Socket, e magic.Empty) magic.AppliedView {
	magic.Assign(s, "child", func() magic.AppliedView {
		time.Sleep(time.Second * 5)
		return test(s, e)
	})
	return homeView(s)
})
