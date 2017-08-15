package chats

import (
	"fmt"
	"time"

	"gopkg.in/vmihailenco/msgpack.v2"

	"github.com/tarantool/go-tarantool"
)

var _ StateStore = (*TarantoolState)(nil)

func NewTarantoolState(conn *tarantool.Connection) *TarantoolState {
	return &TarantoolState{
		Conn: conn,
	}
}

type TarantoolState struct {
	Conn *tarantool.Connection
}

func (s *TarantoolState) Find(chatID string) (*State, error) {
	var tuples []State
	err := s.Conn.SelectTyped(
		"ff_tgbot_statechats",
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
		return nil, ErrNotFound
	}
	return &tuples[0], nil
}

func (s *TarantoolState) Update(obj *State) error {
	obj.UpdatedAt = time.Now()
	_, err := s.Conn.Replace(
		"ff_tgbot_statechats",
		obj,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *State) EncodeMsgpack(e *msgpack.Encoder) error {
	if err := e.EncodeSliceLen(5); err != nil {
		return err
	}

	if err := e.EncodeString(m.ChatID); err != nil {
		return err
	}
	if err := e.EncodeString(m.ScriptID); err != nil {
		return err
	}
	if err := e.EncodeString(m.LastQID); err != nil {
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

func (m *State) DecodeMsgpack(d *msgpack.Decoder) error {
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
	if m.ScriptID, err = d.DecodeString(); err != nil {
		return err
	}
	if m.LastQID, err = d.DecodeString(); err != nil {
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
