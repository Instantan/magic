package magic

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/Instantan/magic/internal"
	"github.com/valyala/fasttemplate"
)

type Template struct {
	filename string

	live   []byte
	static []byte
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
		filename: filename,
		static:   b,
	}
	template.prepareLiveFromStatic()
	return template, nil
}

// MustParseFile is just magic.Must(magic.ParseFile(filename))
func MustParseFile(filename string) *Template {
	return Must(ParseFile(filename))
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
		static: buffer.Bytes(),
	}
	return template, nil
}

// func (templates *Templates) Template(filename string) *Template {
// 	return nil
// }

// Clone copies the template and returns a reference to the new template
func (template *Template) Clone() *Template {
	return &Template{
		filename: template.filename,
		static:   bytes.Clone(template.static),
		live:     bytes.Clone(template.live),
	}
}

// Apply updates the template with all the given placeholders replaced
func (template *Template) Apply(data any) {
	m := dataToMapAny(data)
	template.static = []byte(fasttemplate.ExecuteFuncString(string(template.static), "{{", "}}", func(w io.Writer, tag string) (int, error) {
		if v := jsonGetPath(m, tag); v != nil {
			return w.Write([]byte(fmt.Sprint(v)))
		}
		return w.Write([]byte("{{" + tag + "}}"))
	}))
	template.prepareLiveFromStatic()
}

// Execute writes template with all the given placeholders replaced to the given writer
func (template *Template) ExecuteStatic(w io.Writer, data any) {
	m := dataToMapAny(data)
	fasttemplate.ExecuteFunc(string(template.static), "{{", "}}", w, func(w io.Writer, tag string) (int, error) {
		if v := jsonGetPath(m, tag); v != nil {
			return w.Write([]byte(fmt.Sprint(v)))
		}
		return w.Write([]byte(""))
	})
}

// This should get called after every change to the static data. it turns the template into a live template
func (template *Template) prepareLiveFromStatic() {
	template.live = []byte(internal.ReplaceTemplateBracesInHTMLInnerTextWithComponent(string(template.static)))
}

// Execute writes template with all the given placeholders replaced to the given writer
func (template *Template) executeLiveTemplate(w io.Writer, connId string, data any) {
	b := dataToJSONBytes(data)
	injected := internal.InjectDataIntoHTML(template.live, func() []byte {
		dataSRR := "data-ss=\"" + base64.StdEncoding.EncodeToString(b) + "\""
		dataConnID := "data-connid=\"" + connId + "\""
		return []byte(" " + dataSRR + " " + dataConnID)
	}, injectReactivityScript)
	fasttemplate.ExecuteFunc(string(injected), "{{", "}}", w, func(w io.Writer, tag string) (int, error) {
		return w.Write([]byte("{{" + tag + "}}"))
	})
}

// Writes the data to the given writer
func (template *Template) Write(w io.Writer) {
	w.Write(template.static)
}
