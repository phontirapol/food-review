package template

import (
	"html/template"
	"net/http"
)

type Templater interface {
	ExecuteTemplate(w http.ResponseWriter, template string, data interface{}) error
}

type Template struct {
	Tmpl *template.Template
}

func (t *Template) LoadTemplates(pattern string) error {
	if templates, err := template.ParseGlob(pattern); err != nil {
		return err
	} else {
		t.Tmpl = templates
		return nil
	}
}

func (t *Template) ExecuteTemplate(w http.ResponseWriter, template string, data interface{}) error {
	err := t.Tmpl.ExecuteTemplate(w, template, data)
	if err != nil {
		return err
	}
	return nil
}
