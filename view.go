package magic

import (
	"encoding/json"
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

func (av AppliedView) HTML(w io.Writer) (n int, err error) {
	av.template.Execute(w, func(w io.Writer, tag string) (int, error) {
		if tag == "magic:live" {
			return w.Write(magicMinScript)
		}
		rv, ok := av.socketref.state[tag]
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

func (av AppliedView) Patch() []*patch {
	psBySocketRef := map[uintptr]*patch{}
	av.patch(&psBySocketRef)
	ps := make([]*patch, len(psBySocketRef)+1)
	ps[0] = getPatch()
	ps[0].socketid = socketid(av.socketref.root.id())
	i := 1
	for k := range psBySocketRef {
		ps[i] = psBySocketRef[k]
		i++
	}
	ps[0].data = map[string]any{
		"#": AppliedView{
			socketref: av.socketref,
			template:  av.template,
		},
	}
	return ps
}

func (av AppliedView) patch(patchesBySocketRef *map[uintptr]*patch) {
	refid := av.socketref.id()
	if _, ok := (*patchesBySocketRef)[refid]; ok {
		return
	}
	p := getPatch()
	p.data = av.socketref.state
	p.socketid = socketid(refid)
	(*patchesBySocketRef)[refid] = p
	for _, d := range p.data {
		switch v := d.(type) {
		case AppliedView:
			v.patch(patchesBySocketRef)
		}
	}
}

func (av AppliedView) MarshalJSON() ([]byte, error) {
	d := make([]json.RawMessage, 2)
	d[0] = socketid(av.socketref.id())
	d[1], _ = json.Marshal(av.template.ID())
	return json.Marshal(d)
}
