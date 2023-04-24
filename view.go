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

// TODO: This is unclean af
// reimplement
func (av AppliedView) Patch(templates *[]*Template, knownSocketRefs *map[uintptr]struct{}) []*patch {
	ps := []*patch{}
	p := getPatch()
	rootid, refid := av.socketref.id()
	p.socketid = socketid(rootid, refid)
	p.data = av.socketref.state
	ps = append(ps, p)
	(*knownSocketRefs)[refid] = struct{}{}
	for _, d := range p.data {
		switch v := d.(type) {
		case AppliedView:
			_, irefid := v.socketref.id()
			if _, ok := (*knownSocketRefs)[irefid]; !ok {
				ps = append(ps, v.Patch(knownSocketRefs)...)
			}
		}
	}
	return ps
}

// func (av AppliedView)

// func (av AppliedView) Mount() {
// 	if av.socketref.eventHandler == nil {
// 		return
// 	}
// 	av.socketref.eventHandler(MountEvent, nil)
// }
