package chats

import "errors"

var (
	ErrNotFoundQuestion = errors.New("not found question")
	ErrNotFoundScript   = errors.New("not found script")
	ErrNotFound         = errors.New("not found")
)
