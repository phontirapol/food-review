package template

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

func LoadTemplates(pattern string) error {
	if templates, err := template.ParseGlob(pattern); err != nil {
		return err
	} else {
		tmpl = templates
		return nil
	}
}

func ExecuteTemplate(w http.ResponseWriter, template string, data interface{}) error {
	err := tmpl.ExecuteTemplate(w, template, data)
	if err != nil {
		return err
	}
	return nil
}
