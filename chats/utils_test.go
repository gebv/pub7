package chats

import (
	"testing"
)

func TestTemplateFormat(t *testing.T) {
	type args struct {
		tpl  string
		args map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"",
			args{
				"{{.name}} - {{.ok}}",
				map[string]interface{}{
					"ok": "OKOKOKOK",
				},
			},
			" - OKOKOKOK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format2String(tt.args.tpl, tt.args.args); got != tt.want {
				t.Errorf("TemplateFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
