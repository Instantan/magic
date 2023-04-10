package magic

import (
	"bytes"
	"text/template"
)

type MagicTemplate struct {
	viewResolver func(name string) string
	template     *template.Template
}

func CompileTemplate(tmplt string) *MagicTemplate {
	mt := &MagicTemplate{}
	t, err := template.New("").Funcs(template.FuncMap{}).Parse(tmplt)
	if err != nil {
		panic(err)
	}

	mt.template = t

	return mt
}

func (t MagicTemplate) Exec(data any) string {
	w := bytes.NewBufferString("")
	t.template.Execute(w, data)
	return w.String()
}
