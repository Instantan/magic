package magic

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/Instantan/template"
)

type Template = template.Template

type ViewFn func(s Socket) AppliedView

type AppliedView struct {
	ref      *ref
	template *Template
}

type Views = []AppliedView

type htmlRenderConfig struct {
	magicScriptInline bool
	magicScriptUrl    string
}

func View(templ string) ViewFn {
	t := template.Parse(injectLiveScript(templ))
	return func(s Socket) AppliedView {
		return AppliedView{
			ref:      s.(*ref),
			template: t,
		}
	}
}

func (av AppliedView) html(w io.Writer, config *htmlRenderConfig) (n int, err error) {
	av.template.Execute(w, func(w io.Writer, tag string) (int, error) {
		if tag == "magic:live" {
			if config != nil && !config.magicScriptInline && config.magicScriptUrl != "" {
				return w.Write(unsafeStringToBytes(`<script src=\"` + config.magicScriptUrl + `\" defer/>`))
			}
			return w.Write(magicMinScriptWithTags)
		}
		rv, ok := av.ref.state[tag]
		if !ok {
			return 0, nil
		}
		switch v := rv.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64, bool:
			return w.Write(unsafeStringToBytes(fmt.Sprint(v)))
		case string:
			return w.Write(unsafeStringToBytes(v))
		case AppliedView:
			return v.html(w, config)
		case []AppliedView:
			n := 0
			for i := range v {
				m, err := v[i].html(w, config)
				if err != nil {
					log.Println(err)
				}
				n += m
			}
			return n, nil
		default:
			return 0, nil
		}
	})
	return n, err
}

func (av AppliedView) Patch() []*patch {
	psByref := map[uintptr]*patch{}
	av.patch(&psByref)
	ps := make([]*patch, len(psByref)+1)
	ps[0] = getPatch()
	ps[0].socketid = socketid(av.ref.root.id())
	i := 1
	for k := range psByref {
		ps[i] = psByref[k]
		i++
	}
	ps[0].data = map[string]any{
		"#": AppliedView{
			ref:      av.ref,
			template: av.template,
		},
	}
	return ps
}

func (av AppliedView) patch(patchesByref *map[uintptr]*patch) {
	refid := av.ref.id()
	if _, ok := (*patchesByref)[refid]; ok {
		return
	}
	p := getPatch()
	p.data = av.ref.state
	p.socketid = socketid(refid)
	(*patchesByref)[refid] = p
	for _, d := range p.data {
		switch v := d.(type) {
		case AppliedView:
			v.patch(patchesByref)
		case []AppliedView:
			for i := range v {
				v[i].patch(patchesByref)
			}
		}
	}
}

func (av AppliedView) MarshalJSON() ([]byte, error) {
	d := make([]json.RawMessage, 2)
	d[0] = socketid(av.ref.id())
	d[1], _ = json.Marshal(av.template.ID())
	return json.Marshal(d)
}

func (av AppliedView) marshalPatchJSON() ([]byte, error) {
	m := make([]json.RawMessage, 2)
	m[0], _ = json.Marshal(av.template.ID())
	m[1], _ = json.Marshal(av.template.String())
	return json.Marshal(m)
}
