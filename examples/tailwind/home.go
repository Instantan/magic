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
		<link href="/assets/index.min.css" rel="stylesheet">
	</head>

	<body class="bg-black text-whited">
		Test!
	</body>
	</html>
`)

var home = magic.Component(func(s magic.Socket, _ magic.Empty) magic.AppliedView {
	return homeView(s)
})
