package magic

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"

	"github.com/valyala/fasttemplate"
)

type Template struct {
	data []byte
}

type Templates []*Template

// // ParseFiles reads the given files from the filepath and returns them as a slice of templates
// func ParseFiles(filenames ...string) (Templates, error) {
// 	templates := Templates{}
// 	for _, filename := range filenames {
// 		template, err := ParseFile(filename)
// 		if err != nil {
// 			return nil, err
// 		}
// 		templates = append(templates, template)
// 	}
// 	return templates, nil
// }

// ParseFile reads the given file from the filepath and returns it as a template
func ParseFile(filename string) (*Template, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	template := &Template{
		data: b,
	}
	return template, nil
}

// // ParseFiles reads the given files from the filesystem and returns them as a slice of templates
// func ParseFilesFS(fsys fs.FS, filenames ...string) (Templates, error) {
// 	templates := Templates{}
// 	for _, filename := range filenames {
// 		template, err := ParseFileFS(fsys, filename)
// 		if err != nil {
// 			return nil, err
// 		}
// 		templates = append(templates, template)
// 	}
// 	return templates, nil
// }

// ParseFileFS reads the given file from the given filesystem and returns it as a template
func ParseFileFS(fsys fs.FS, filename string) (*Template, error) {
	file, err := fsys.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buffer := bytes.Buffer{}
	_, err = buffer.ReadFrom(file)
	if !errors.Is(err, io.EOF) {
		return nil, err
	}
	template := &Template{
		data: buffer.Bytes(),
	}
	return template, nil
}

// func (templates *Templates) Template(filename string) *Template {
// 	return nil
// }

// Clone copies the template and returns a reference to the new template
func (template *Template) Clone() *Template {
	return &Template{
		data: bytes.Clone(template.data),
	}
}

// Apply updates the template with all the given placeholders replaced
func (template *Template) Apply(data map[string]string) {
	template.data = []byte(fasttemplate.ExecuteFuncString(string(template.data), "{{", "}}", func(w io.Writer, tag string) (int, error) {
		if v, ok := data[tag]; ok {
			return w.Write([]byte(v))
		}
		return w.Write([]byte("{{" + tag + "}}"))
	}))
}

// Execute writes template with all the given placeholders replaced to the given writer
func (template *Template) Execute(w io.Writer, data map[string]string) {
	fasttemplate.ExecuteFunc(string(template.data), "{{", "}}", w, func(w io.Writer, tag string) (int, error) {
		if v, ok := data[tag]; ok {
			return w.Write([]byte(v))
		}
		return w.Write([]byte("{{" + tag + "}}"))
	})
}

// Writes the data to the given writer
func (template *Template) Write(w io.Writer) {
	w.Write(template.data)
}
