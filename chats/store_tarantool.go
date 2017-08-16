package chats

import (
	"fmt"
	"time"

	"gopkg.in/vmihailenco/msgpack.v2"

	"github.com/gebv/ff_tgbot/errors"
	"github.com/tarantool/go-tarantool"
)

var _ Store = (*TarantoolChat)(nil)

func NewTarantool(conn *tarantool.Connection) *TarantoolChat {
	return &TarantoolChat{
		Conn: conn,
	}
}

type TarantoolChat struct {
	Conn *tarantool.Connection
}

func (s *TarantoolChat) Find(chatID string) (*Chat, error) {
	var tuples []Chat
	err := s.Conn.SelectTyped(
		"ff_chats",
		"primary",
		0,
		1,
		tarantool.IterEq,
		[]interface{}{
			chatID,
		},
		&tuples,
	)
	if err != nil {
		return nil, err
	}
	if len(tuples) == 0 {
		return nil, errors.ErrNotFound
	}
	return &tuples[0], nil
}

func (s *TarantoolChat) Update(obj *Chat) error {
	obj.UpdatedAt = time.Now()
	_, err := s.Conn.Replace(
		"ff_chats",
		obj,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *Chat) EncodeMsgpack(e *msgpack.Encoder) error {
	if err := e.EncodeSliceLen(5); err != nil {
		return err
	}

	if err := e.EncodeString(m.ChatID); err != nil {
		return err
	}
	if err := e.EncodeString(m.PreviousNodeID); err != nil {
		return err
	}
	if err := e.EncodeString(m.NextNodeID); err != nil {
		return err
	}
	if err := e.Encode(m.Props); err != nil {
		return err
	}
	if err := e.EncodeTime(m.UpdatedAt); err != nil {
		return err
	}
	return nil
}

func (m *Chat) DecodeMsgpack(d *msgpack.Decoder) error {
	var err error
	var l int
	if l, err = d.DecodeSliceLen(); err != nil {
		return err
	}
	if l != 5 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}

	if m.ChatID, err = d.DecodeString(); err != nil {
		return err
	}
	if m.PreviousNodeID, err = d.DecodeString(); err != nil {
		return err
	}
	if m.NextNodeID, err = d.DecodeString(); err != nil {
		return err
	}
	if err = d.Decode(&m.Props); err != nil {
		return err
	}
	if m.UpdatedAt, err = d.DecodeTime(); err != nil {
		return err
	}
	return nil
}
