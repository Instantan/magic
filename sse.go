package magic

import (
	"net/http"
)

type sse struct {
	w http.ResponseWriter
}

func establishSSEConnection(w http.ResponseWriter) *sse {
	_, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return nil
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	return &sse{
		w: w,
	}
}

func (e *sse) Write(p []byte) (n int, err error) {
	p = append(append([]byte("data: "), p...), []byte("\n\n")...)
	n, err = e.w.Write(p)
	e.w.(http.Flusher).Flush()
	return n, err
}
