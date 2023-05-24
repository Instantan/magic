package magic

type AnimateArgs struct {
	Keyframes []AnimationKeyframe `json:"keyframes"`
	Options   AnimateOptions      `json:"options,omitempty"`
	Duration  int                 `json:"duration,omitempty"`
}

type AnimationKeyframe map[string]string

type AnimateOptions struct {
	Delay              int     `json:"delay,omitempty"`
	Direction          string  `json:"direction,omitempty"`
	Duration           int     `json:"duration,omitempty"`
	Easing             string  `json:"easing,omitempty"`
	EndDelay           int     `json:"endDelay,omitempty"`
	Fill               string  `json:"fill,omitempty"`
	IterationStart     float64 `json:"iterationStart,omitempty"`
	Iterations         string  `json:"iterations,omitempty"`
	Composite          string  `json:"composite,omitempty"`
	IterationComposite string  `json:"iterationComposite,omitempty"`
	PseudoElement      string  `json:"pseudoElement,omitempty"`
}

type ScrollIntoViewArgs struct {
	AlignToTop bool                  `json:"alignToTop,omitempty"`
	Options    ScrollIntoViewOptions `json:"options,omitempty"`
}

type ScrollIntoViewOptions struct {
	Behavior string `json:"behavior,omitempty"`
	Block    string `json:"block,omitempty"`
	Inline   string `json:"inline,omitempty"`
}

// Sends the reload event to the given socket, the reload is not a full reload, instead its a websocket reconnect
func Reload(s Socket) {
	s.DispatchEvent(NavigateEvent, s.Request().URL.String())
}

// Sends a navigation event to the given socket, the navigate is a live navigation
func Navigate(s Socket, location string) {
	s.DispatchEvent(NavigateEvent, location)
}

// Calls the reset() function for the given ids (if no id is given it uses the id of the socket)
func Reset(s Socket, ids ...string) {
	s.DispatchEvent(ResetEvent, ids)
}

// Calls the click() function for the given ids (if no id is given it uses the id of the socket)
func Click(s Socket, ids ...string) {
	s.DispatchEvent(ClickEvent, ids)
}

// Calls the blur() function for the given ids (if no id is given it uses the id of the socket)
func Blur(s Socket, ids ...string) {
	s.DispatchEvent(ClickEvent, ids)
}

// Calls the focus() function for the given ids (if no id is given it uses the id of the socket)
func Focus(s Socket, ids ...string) {
	s.DispatchEvent(FocusEvent, ids)
}

// Enters the fullscreen mode of the window
func OpenFullscreen(s Socket) {
	s.DispatchEvent(OpenFullscreenEvent, nil)
}

// Leaves the fullscreen mode of the window
func CloseFullscreen(s Socket) {
	s.DispatchEvent(CloseFullscreenEvent, nil)
}

// Calls the animate(keyframes, options) function for the given ids (if no id is given it uses the id of the socket)
func Animate(s Socket, args AnimateArgs, ids ...string) {
	s.DispatchEvent(AnimateEvent, struct {
		Args AnimateArgs `json:"args"`
		Ids  []string    `json:"ids"`
	}{args, ids})
}

// Calls the scrollIntoView(options) function for the given ids (if no id is given it uses the id of the socket)
func ScrollIntoView(s Socket, args ScrollIntoViewArgs, ids ...string) {
	s.DispatchEvent(AnimateEvent, struct {
		Args ScrollIntoViewArgs `json:"args"`
		Ids  []string           `json:"ids"`
	}{args, ids})
}

// Disconnects the socket from the client side. It wont reconnect afterwards
func Disconnect(s Socket) {
	s.DispatchEvent(DisconnectEvent, nil)
}

// UpdateURL changes the clients url to the request url without navigating
func UpdateURL(s Socket) {
	s.DispatchEvent(UpdateURLEvent, urlToStringWithoutSchemeAndHost(s.Request().URL))
}
