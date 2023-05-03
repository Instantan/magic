package main

import (
	"time"

	"github.com/Instantan/magic"
)

var homeView = magic.View(`
	<!DOCTYPE html>
	<html lang="de">

	<head>
		<meta charset="UTF-8">
		<title>Test</title>
	</head>
	<style>
		html {
			background-color: black;
			color: rgb(200, 200, 200);
		}
		
		.true {
			color: blueviolet;
		}
	</style>

	<body class="body">
		{{navbar}}
		{{liveNavbar}}
	</body>
	</html>
`)

var home = magic.Component(func(s magic.Socket) magic.AppliedView {
	magic.Assign(s, "navbar", navbarComponent(s))
	if s.Live() {
		t := time.NewTicker(time.Second * 5)
		go func() {
			for {
				select {
				case <-t.C:
					magic.Assign(s, "liveNavbar", counterComponent(s))
				case <-s.Done():
					t.Stop()
					return
				}
			}
		}()
	}
	return homeView(s)
})
