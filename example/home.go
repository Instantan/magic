package main

import (
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
	</body>
	</html>
`)

var home = magic.Component(func(s magic.Socket) magic.AppliedView {
	// i := 0
	// s.State().Set("navbar", navbarView(s))
	magic.Assign(s, "navbar", navbarView(s))
	// if s.Live() {
	// 	go func() {
	// 		for {
	// 			select {
	// 			case <-s.Done():
	// 				return
	// 			default:
	// 				time.Sleep(time.Second)
	// 				i++
	// 				println(i)
	// 				magic.Assign(s, "navbar", i)
	// 			}
	// 		}
	// 	}()
	// }
	return homeView(s)
})
