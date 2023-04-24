package main

import "github.com/Instantan/magic"

var homeView = magic.View(`
	<!DOCTYPE html>
	<html lang="de">

	<head>
		<meta charset="UTF-8">
		<title>Test</title>
	</head>
	<style>
		html {
			background-color: rgb(30, 30, 30);
			color: rgb(200, 200, 200);
		}
		
		.true {
			color: blueviolet;
		}
	</style>

	<body class="body">
		{{navbar}}
	</body>
	</html>
`)

var home = magic.Component(func(s magic.Socket) magic.AppliedView {
	// s.State().Set("navbar", navbarView(s))

	magic.Assign(s, "navbar", navbarView(s))

	return homeView(s)
})
