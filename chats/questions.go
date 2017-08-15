package chats

import (
	lua "github.com/yuin/gopher-lua"
)

type Question struct {
	QID string `toml:"id"`

	// question text
	Texts []string `toml:"texts"`

	LuaScript string `toml:"script"`

	// сохранить ответ как параметр
	ParamName string `toml:"param_name"`

	Options []OptionOfResponse `toml:"opts"`

	// глобальное значение следующего вопроса
	// либо будет переопределен из Options
	NextQID string `toml:"next_qid"`
}

func (q Question) LastText() string {
	return q.Texts[len(q.Texts)-1]
}

func (q Question) TextsWithoutLast() []string {
	if len(q.Texts) == 0 {
		return []string{}
	}
	return q.Texts[0 : len(q.Texts)-1]
}

func (q Question) ExecuteScript(props map[string]interface{}) error {
	if len(q.LuaScript) == 0 {
		return nil
	}
	L := lua.NewState()
	defer L.Close()
	registerContextType(L, &LuaContext{Props: props})
	return L.DoString(q.LuaScript)
}

type OptionOfResponse struct {
	Key     string `toml:"key"`
	NextQID string `toml:"qid"`
}
