package main

import (
	"github.com/Instantan/magic"
)

type HTMLProps struct {
	title string
	body  any
}

var htmlView = magic.View(`
	<!DOCTYPE html>
	<html lang="de">

	<head>
		<meta charset="UTF-8">
		<meta name="theme-color" content="#35B6D2" />
		<title>{{title}}</title>
	</head>
	<style>
		html {
			background-color: black;
			color: rgb(200, 200, 200);
		}
		
		.true {
			color: blueviolet;
		}
		.magic-connecting {
			background-color: white;
		}
	</style>

	<body class="body">
		{{body}}
	</body>
	</html>
`)

var html = magic.Component(func(s magic.Socket, props HTMLProps) magic.AppliedView {
	magic.Assign(s, "title", props.title)
	magic.Assign(s, "body", props.body.(magic.AppliedView))
	return htmlView(s)
})
