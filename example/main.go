package main

import (
	"log"
	"net/http"

	"github.com/Instantan/magic"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/index.html", magic.CreatePage(magic.MustParseFile("index.html"), IndexPage))

	log.Print("Listening to http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
