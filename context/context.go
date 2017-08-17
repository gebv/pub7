package context

import (
	"fmt"

	"github.com/gebv/ff_tgbot/utils"
	lua "github.com/yuin/gopher-lua"
)

const luaCtxTypeName = "ctx"

type Context interface {
	Props() map[string]interface{}
	Set(string, interface{})
	Get(string) interface{}

	IsAbort() bool
	Abort()

	RedirectTo() string
	SetRedirect(id string)

	CurrentTextMessage() string

	Error() error
}

func luaCheckCtx(L *lua.LState) Context {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(Context); ok {
		return v
	}
	reason := fmt.Sprintf("expected Context, got %T", ud.Value)
	L.ArgError(1, reason)
	return nil
}

var basicMethods = map[string]lua.LGFunction{
	"set": func(L *lua.LState) int {
		ctx := luaCheckCtx(L)
		k := L.CheckString(2)
		lv := L.CheckAny(3)
		v := utils.ToValueFromLValue(lv)
		ctx.Set(k, v)
		return 0
	},
	"get": func(L *lua.LState) int {
		ctx := luaCheckCtx(L)
		k := L.CheckString(2)
		v := ctx.Get(k)
		lv := utils.ToLValueOrNil(v, L)
		L.Push(lv)
		return 1
	},
	"text": func(L *lua.LState) int {
		ctx := luaCheckCtx(L)
		lv := utils.ToLValueOrNil(ctx.CurrentTextMessage(), L)
		L.Push(lv)
		return 1
	},
	"abort": func(L *lua.LState) int {
		ctx := luaCheckCtx(L)
		ctx.Abort()
		return 0
	},
	"redirect": func(L *lua.LState) int {
		ctx := luaCheckCtx(L)
		ctx.SetRedirect(L.CheckString(2))
		return 0
	},
}
