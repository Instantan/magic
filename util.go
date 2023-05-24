package magic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"unsafe"
)

type Empty = struct{}
type Nothing = Empty

// Must panics if the passed error is not nil
func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

func Map[T any, R any](s []T, f func(e T) R) []R {
	ns := make([]R, len(s))
	for i := range s {
		ns[i] = f(s[i])
	}
	return ns
}

func socketid(id uintptr) json.RawMessage {
	v, _ := json.Marshal(fmt.Sprintf("%v", id))
	return v
}

func unsafeStringToBytes(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

func urlToStringWithoutSchemeAndHost(u *url.URL) string {
	var buf strings.Builder
	if u.Opaque != "" {
		buf.WriteString(u.Opaque)
	} else {
		path := u.EscapedPath()
		if path != "" && path[0] != '/' && u.Host != "" {
			buf.WriteByte('/')
		}
		if buf.Len() == 0 {
			// RFC 3986 §4.2
			// A path segment that contains a colon character (e.g., "this:that")
			// cannot be used as the first segment of a relative-path reference, as
			// it would be mistaken for a scheme name. Such a segment must be
			// preceded by a dot-segment (e.g., "./this:that") to make a relative-
			// path reference.
			if segment, _, _ := strings.Cut(path, "/"); strings.Contains(segment, ":") {
				buf.WriteString("./")
			}
		}
		buf.WriteString(path)
	}
	if u.ForceQuery || u.RawQuery != "" {
		buf.WriteByte('?')
		buf.WriteString(u.RawQuery)
	}
	if u.Fragment != "" {
		buf.WriteByte('#')
		buf.WriteString(u.EscapedFragment())
	}
	return buf.String()
}

//go:linkname acceptsGzip gzhttp.acceptsGzip
func acceptsGzip(r *http.Request) bool
