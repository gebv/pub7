package utils

import (
	"bytes"
	"html/template"
)

func ExecuteTemplate(tpl string, args map[string]interface{}) string {
	t := template.Must(template.New("t1").Parse(tpl))
	var buf bytes.Buffer
	t.Execute(&buf, args)
	return buf.String()
}
