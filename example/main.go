package main

import (
	"log"
	"net/http"

	"github.com/Instantan/magic"
)

func main() {
	mux := http.NewServeMux()

	template := magic.Must(magic.Tem)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
