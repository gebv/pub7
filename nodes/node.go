package nodes

import (
	"github.com/gebv/ff_tgbot/context"

	lua "github.com/yuin/gopher-lua"
)

type Store interface {
	Find(chatID string) (*Node, error)
	LoadFromToml(dat []byte) error
}

type Node struct {
	ID   string `toml:"id"`
	Text Texts  `toml:"text"`

	BeforeScript  string `toml:"before"`
	HandlerString string `toml:"handler"`
	AfterScript   string `toml:"after"`

	NextNodeID string `toml:"next"`
	ParamName  string `toml:"param"`
	IsTransit  bool   `toml:"transit"`

	Options []NodeOption `toml:"opts"`
}

func (n *Node) Before(L *lua.LState, ctx context.Context) (err error) {
	return n.executeLuaScript(L, n.BeforeScript, ctx)
}

func (n *Node) Handler(L *lua.LState, ctx context.Context) error {
	return n.executeLuaScript(L, n.HandlerString, ctx)
}

func (n *Node) After(L *lua.LState, ctx context.Context) error {
	return n.executeLuaScript(L, n.AfterScript, ctx)
}

func (n *Node) executeLuaScript(L *lua.LState, script string, ctx context.Context) error {
	return L.DoString(script)
}

type NodeOption struct {
	Text       string `toml:"text"`
	NextNodeID string `toml:"next"`
}
