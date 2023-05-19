package main

import (
	"encoding/json"
	"fmt"

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
	</style>

	<body class="body">
		<form action="/" magic:submit="/">
			<label for="name">Name:</label><br>
			<input type="text" id="name" name="name" />
			<input type="submit" value="Submit" />
		</form>
	</body>
	</html>
`)

var home = magic.Component(func(s magic.Socket, e magic.Empty) magic.AppliedView {
	s.HandleEvent(func(ev string, data magic.EventData) {
		switch ev {
		case magic.SubmitEvent:
			s := data.Submit()
			m := make(map[string]any)
			json.Unmarshal(s.Form, &m)
			fmt.Printf("%v - %v", ev, m)
		}
	})
	return homeView(s)
})
