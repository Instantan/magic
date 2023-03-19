package magic

import (
	"context"
	"net/http"
)

type LiveTemplate struct {
}

type EventValue struct{}
type Socket[T any] struct {
	Context context.Context

	State T // this State can be updated

	pstate []byte // the state the client holds, gets updated after every patch
}

// type LiveHandler[T any] interface {
// 	Mount(socket *Socket[T])
// 	EventListener(event string, value EventValue, socket *Socket[T])
// 	Unmount(socket *Socket[T])
// }

// Live turns the template into a live template (copies the underlying template)
func (template *Template) Live(handler LiveTemplate) *LiveTemplate {
	for i := 0; i < len(template.data); i++ {

	}
	return nil
}

// Execute sends the live template to the writer with the given data
func (liveTemplate *LiveTemplate) Execute(w http.ResponseWriter, data map[string]string) {

}

// DispatchEvent sends a event to the live client
func (sock *Socket[T]) DispatchEvent(event string, value EventValue) {

}

// DispatchPatch sends a json patch to the live client
// the patch is calculated based on the changes of the
func (sock *Socket[T]) DispatchPatch() {

}
