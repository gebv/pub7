package chats

import (
	"time"
)

type StateStore interface {
	Find(chatID string) (*State, error)
	Update(obj *State) error
}

func NewState(chatID string) *State {
	return &State{
		ChatID: chatID,
		Props:  make(map[string]interface{}),
	}
}

type State struct {
	ChatID   string
	ScriptID string
	LastQID  string

	Props map[string]interface{}

	UpdatedAt time.Time
}

func (s State) Prop(name string) interface{} {
	return s.Props[name]
}
