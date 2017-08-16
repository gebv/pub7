package context

import (
	lua "github.com/yuin/gopher-lua"
	"gopkg.in/telegram-bot-api.v4"
)

var _ Context = (*TelegramContext)(nil)

func Telegram(
	props map[string]interface{},
	api *tgbotapi.BotAPI,
) *TelegramContext {
	return &TelegramContext{
		props: props,
		api:   api,
	}
}

type TelegramContext struct {
	props   map[string]interface{}
	isAbort bool
	err     error
	redirTo string

	chatID int64
	api    *tgbotapi.BotAPI
}

func RegisterTelegramContextType(
	L *lua.LState,
	ctx *TelegramContext,
) {
	mt := L.NewTypeMetatable(luaCtxTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), joinFuncs(basicMethods, telegramMethods)))

	ud := L.NewUserData()
	ud.Value = ctx
	L.SetMetatable(ud, L.GetTypeMetatable(luaCtxTypeName))

	L.SetGlobal("ctx", ud)
	return
}

var telegramMethods = map[string]lua.LGFunction{}

func joinFuncs(m1, m2 map[string]lua.LGFunction) map[string]lua.LGFunction {
	list := make(map[string]lua.LGFunction, len(m1)+len(m2))
	for k, fn := range m1 {
		list[k] = fn
	}
	for k, fn := range m2 {
		list[k] = fn
	}
	return list
}

// implement the Context

func (c *TelegramContext) Props() map[string]interface{} {
	return c.props
}

func (c *TelegramContext) Set(k string, v interface{}) {
	c.props[k] = v
}

func (c *TelegramContext) Get(k string) interface{} {
	return c.props[k]
}

func (c *TelegramContext) IsAbort() bool {
	return c.isAbort
}

func (c *TelegramContext) Abort() {
	c.isAbort = true
}

func (c *TelegramContext) RedirectTo() string {
	return c.redirTo
}

func (c *TelegramContext) SetRedirect(v string) {
	c.redirTo = v
}

func (c *TelegramContext) Error() error {
	return c.err
}
