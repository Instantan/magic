package main

import (
	"github.com/Instantan/magic"
)

var homeView = magic.View(`
	<!DOCTYPE html>
	<html lang="de">

	<head>
		<meta charset="UTF-8">
		<meta name="theme-color" content="#35B6D2" />
		<title>Test</title>
	</head>

	<body>
	</body>
	</html>
`)

var home = magic.Component(func(s magic.Socket) magic.AppliedView {
	return homeView(s)
})
