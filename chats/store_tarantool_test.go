package chats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/vmihailenco/msgpack.v2"
)

func TestChatMsgpackSerializer(t *testing.T) {
	want := &Chat{
		ChatID:         "chats:1",
		PreviousNodeID: "scripts:1",
		NextNodeID:     "questions:1",
		Props: map[string]interface{}{
			"a":   "b",
			"c":   "d",
			"int": float64(10),
		},
		UpdatedAt: time.Now(),
	}

	dat, err := msgpack.Marshal(want)
	assert.NoError(t, err, "marshal Chat")

	got := &Chat{}

	err = msgpack.Unmarshal(dat, got)
	assert.NoError(t, err, "unmarshal Chat")

	assert.Equal(t, want, got)
}
