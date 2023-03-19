package magic

import "net/http"

// PageRenderer is a simple function that returns a structure
// the returned structure is used to render the template
type PageRenderer func(context PageContext) any

type PageContext struct {
	Request *http.Request
}

type Page struct {
}

func CreatePage(renderer PageRenderer) *Page {
	return &Page{}
}
