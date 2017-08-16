package chats

import (
	"time"
)

type Store interface {
	Find(chatID string) (*Chat, error)
	Update(obj *Chat) error
}

func NewChat(chatID string) *Chat {
	return &Chat{
		ChatID: chatID,
		Props:  make(map[string]interface{}),
	}
}

type Chat struct {
	ChatID         string
	PreviousNodeID string
	NextNodeID     string
	Props          map[string]interface{}
	UpdatedAt      time.Time
}

func (s Chat) Prop(name string) interface{} {
	return s.Props[name]
}
