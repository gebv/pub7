package chats

import (
	"time"

	"github.com/gebv/ff_tgbot/errors"
)

var _ Store = (*InmemoryChat)(nil)

func NewInmemory() *InmemoryChat {
	return &InmemoryChat{
		List: make(map[string]Chat),
	}
}

type InmemoryChat struct {
	List map[string]Chat
}

func (s *InmemoryChat) Find(chatID string) (*Chat, error) {
	obj, exists := s.List[chatID]
	if !exists {
		return nil, errors.ErrNotFound
	}
	return &obj, nil
}

func (s *InmemoryChat) Update(obj *Chat) error {
	obj.UpdatedAt = time.Now()
	s.List[obj.ChatID] = *obj
	return nil
}
