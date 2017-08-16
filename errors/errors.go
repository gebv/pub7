package errors

import "errors"

var (
	ErrNotFound = errors.New("not found")
)

func New(msg string) error {
	return errors.New(msg)
}
