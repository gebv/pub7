package chats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/vmihailenco/msgpack.v2"
)

func TestStateMsgpackSerializer(t *testing.T) {
	want := &State{
		ChatID:   "chats:1",
		ScriptID: "scripts:1",
		LastQID:  "questions:1",
		Props: map[string]interface{}{
			"a":   "b",
			"c":   "d",
			"int": float64(10),
		},
		UpdatedAt: time.Now(),
	}

	dat, err := msgpack.Marshal(want)
	assert.NoError(t, err, "marshal state")

	got := &State{}

	err = msgpack.Unmarshal(dat, got)
	assert.NoError(t, err, "unmarshal state")

	assert.Equal(t, want, got)
}
