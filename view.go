package magic

import (
	"fmt"
	"io"

	"github.com/Instantan/template"
)

type Template = template.Template

type ViewFn func(s Socket) AppliedView

type AppliedView struct {
	socketref *socketref
	template  *Template
}

func View(templ string) ViewFn {
	t := template.Parse(injectLiveScript(templ))
	return func(s Socket) AppliedView {
		return AppliedView{
			socketref: s.(*socketref),
			template:  t,
		}
	}
}

func (r AppliedView) HTML(w io.Writer) (n int, err error) {
	r.template.Execute(w, func(w io.Writer, tag string) (int, error) {
		if tag == "magic:live" {
			return w.Write(magicMinScript)
		}
		rv, ok := r.socketref.state[tag]
		if !ok {
			return 0, nil
		}
		switch v := rv.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64, bool:
			return w.Write(unsafeStringToBytes(fmt.Sprint(v)))
		case string:
			return w.Write(unsafeStringToBytes(v))
		case AppliedView:
			return v.HTML(w)
		case []AppliedView:
			n := 0
			for i := range v {
				m, _ := v[i].HTML(w)
				m += n
			}
			return n, nil
		default:
			return 0, nil
		}
	})
	return n, err
}
