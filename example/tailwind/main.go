package main

import (
	"log"
	"net/http"

	"github.com/Instantan/magic"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", magic.CompressedComponentHTTPHandler(home))
	mux.Handle("/assets/index.min.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/index.min.css")
	}))

	log.Print("Listening to http://localhost:8070")
	if err := http.ListenAndServe(":8070", mux); err != nil {
		log.Fatal(err)
	}
}
