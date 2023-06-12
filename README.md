# magic

<img src="https://raw.githubusercontent.com/Instantan/magic/main/website/assets/images/magic-logo.svg?token=GHSAT0AAAAAAB7YBL5YY7SBTOQFTZ52KSNIZEHSFYQ" align="right"
     title="Magic logo" width="120">

Magic is a powerful, yet easy-to-use web development framework written in Golang that simplifies the web development process. With its small JavaScript footprint of under 12kb (4.6kb gzipped), Magic helps you build fast and efficient web applications.

# Getting started


For a simple live counter example:

> main.go
```go
package main

func main() {
    mux := http.NewServeMux()
	mux.Handle("/", magic.CompressedComponentHTTPHandler(home))
	log.Print("Listening to http://localhost:8070")
	if err := http.ListenAndServe(":8070", mux); err != nil {
		log.Fatal(err)
	}
}
```

> home.go
```go
package main

import (
	"time"

	"github.com/Instantan/magic"
)

var homeView = magic.View(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Time</title>
	</head>
	<body>
		<p>{{counter}}</p>
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
                t.Stop()
                return
			}
		}
	})
	return homeView(s)
})

```

# Roadmap

- Clean up socket.go, ref.go and view.go
- Benchmarking, memory reduction and performance improvements
- Unit tests