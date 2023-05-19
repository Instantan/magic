package magic

// Sends the reload event to the given socket
// the reload is not a full reload, instead its a websocket reconnect
func Reload(s Socket) {
	s.DispatchEvent(NavigateEvent, s.Request().URL.String())
}

// Sends a navigation event to the given socket
// the navigate is a live navigation
func Navigate(s Socket, location string) {
	s.DispatchEvent(NavigateEvent, location)
}
