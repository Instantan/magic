package magic

import "net/http"

type HandlerFunc[T any] ComponentFn[T]

func (f HandlerFunc[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// f(nil)
}
