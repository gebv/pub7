package chats

import (
	"fmt"

	"github.com/gebv/ff_tgbot/utils"
	lua "github.com/yuin/gopher-lua"
)

const luaCtxTypeName = "ctx"

func NewLuaContext(args map[string]interface{}) *LuaContext {
	return &LuaContext{
		Props: args,
	}
}

type LuaContext struct {
	Props map[string]interface{}
}

func registerContextType(L *lua.LState, c *LuaContext) {
	mt := L.NewTypeMetatable(luaCtxTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), luaCtxMethods))

	L.SetGlobal("ctx", newContextLuaValue(c, L))
	return
}

func newContextLuaValue(ctx *LuaContext, L *lua.LState) lua.LValue {
	ud := L.NewUserData()
	ud.Value = ctx
	L.SetMetatable(ud, L.GetTypeMetatable(luaCtxTypeName))
	return ud
}

func luaCheckCtx(L *lua.LState) *LuaContext {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*LuaContext); ok {
		return v
	}
	reason := fmt.Sprintf("expected Context, got %T", ud.Value)
	L.ArgError(1, reason)
	return nil
}

var luaCtxMethods = map[string]lua.LGFunction{
	"set": func(L *lua.LState) int {
		ctx := luaCheckCtx(L)
		k := L.CheckString(2)
		lv := L.CheckAny(3)
		v := utils.ToValueFromLValue(lv)
		ctx.Props[k] = v
		return 0
	},
	"get": func(L *lua.LState) int {
		ctx := luaCheckCtx(L)
		k := L.CheckString(2)
		v := ctx.Props[k]
		lv := utils.ToLValueOrNil(v, L)
		L.Push(lv)
		return 1
	},
}
