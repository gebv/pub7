package chats

import (
	"bytes"
	"html/template"
)

func Format2String(tpl string, args map[string]interface{}) (res string) {
	t := template.Must(template.New("t1").Parse(tpl))
	var buf bytes.Buffer
	t.Execute(&buf, args)
	return buf.String()
}
