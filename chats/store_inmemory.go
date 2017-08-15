package chats

import (
	"time"
)

var _ StateStore = (*InmemoryState)(nil)

func NewInmemoryState() *InmemoryState {
	return &InmemoryState{
		List: make(map[string]State),
	}
}

type InmemoryState struct {
	List map[string]State
}

func (s *InmemoryState) Find(chatID string) (*State, error) {
	obj, exists := s.List[chatID]
	if !exists {
		return nil, ErrNotFound
	}
	return &obj, nil
}

func (s *InmemoryState) Update(obj *State) error {
	obj.UpdatedAt = time.Now()
	s.List[obj.ChatID] = *obj
	return nil
}
