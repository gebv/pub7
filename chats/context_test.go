package chats

import "testing"
import "github.com/stretchr/testify/assert"

func TestLuaContext(t *testing.T) {
	tests := []struct {
		script string
		ctx    *LuaContext
		want   *LuaContext
	}{
		{
			`
ctx:set("a",ctx:get("d").."qwe")
ctx:set("int",ctx:get("int")+1)
ctx:set("checkInt", ctx:get("int")==6)
`,
			NewLuaContext(map[string]interface{}{"a": "c", "d": "e", "int": "5"}),
			NewLuaContext(map[string]interface{}{"a": "eqwe", "d": "e", "int": float64(6), "checkInt": true}),
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			q := &Question{
				LuaScript: tt.script,
			}
			err := q.ExecuteScript(tt.ctx.Props)
			assert.NoError(t, err, "execute lua")
			assert.Equal(t, tt.want, tt.ctx)
		})
	}
}
