package main

import (
	"log"
	"net/http"

	"github.com/Instantan/magic"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/index.html", magic.HandlerFunc(home))
	mux.Handle("/overview.html", magic.HandlerFunc(overview))

	log.Print("Listening to http://localhost:8070")
	if err := http.ListenAndServe(":8070", mux); err != nil {
		log.Fatal(err)
	}
}
