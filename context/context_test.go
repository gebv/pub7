package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLuaContext(t *testing.T) {
	tests := []struct {
		script string
		ctx    *Context
		want   *Context
	}{
		{
			`
ctx:set("a",ctx:get("d").."qwe")
ctx:set("int",ctx:get("int")+1)
ctx:set("checkInt", ctx:get("int")==6)
`,
			NewContext(map[string]interface{}{"a": "c", "d": "e", "int": "5"}),
			NewContext(map[string]interface{}{"a": "eqwe", "d": "e", "int": float64(6), "checkInt": true}),
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			q := &Node{
				HandlerString: tt.script,
			}
			err := q.Handler(tt.ctx)
			assert.NoError(t, err, "execute lua")
			assert.Equal(t, tt.want, tt.ctx)
		})
	}
}
