package message

import (
	"bytes"
	"strings"
	"text/template"
	"time"
)

var funcMap = template.FuncMap{
	"upper": strings.ToUpper,
	"formatDate": func(t time.Time, layout string) string {
		return t.Format(layout)
	},
}

type Templater struct {
	tmpl *template.Template
}

func NewTemplater(text string) (*Templater, error) {
	t, err := template.New("msg").Funcs(funcMap).Parse(text)
	if err != nil {
		return nil, err
	}
	return &Templater{tmpl: t}, nil
}

func (t *Templater) Render(data map[string]interface{}) (string, error) {
	var buf bytes.Buffer
	if err := t.tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
