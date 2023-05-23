package main

import (
	"fmt"
	"sync/atomic"
	"time"

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

		<input type="text" magic:keypress="test">

		{{navbar}}
		{{liveNavbar}}
	</body>
	</html>
`)

var connectedUsers = &atomic.Int64{}

var home = magic.Component(func(s magic.Socket, e magic.Empty) magic.AppliedView {
	connectedUsers.Add(1)
	fmt.Printf("Connected: %v\n", connectedUsers.Load())
	magic.Assign(s, "navbar", navbarComponent(s, e))
	magic.Assign(s, "liveNavbar", counterComponent(s, e))

	if s.Live() {

		t := time.NewTicker(time.Second * 5)
		go func() {
			for range t.C {
				magic.Assign(s, "liveNavbar", counterComponent(s, e))
			}
		}()
		s.HandleEvent(func(ev string, data magic.EventData) {
			switch ev {
			case magic.KeypressEvent:
				kp := data.Keypress()
				println(kp.Content + kp.Key)
			case magic.UnmountEvent:
				connectedUsers.Add(-1)
				fmt.Printf("Connected: %v\v", connectedUsers.Load())
				t.Stop()
			}
		})
	}
	return homeView(s)
})
