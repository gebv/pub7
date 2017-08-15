package chats

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func TestDecodeWorkspaceFromToml(t *testing.T) {
	want := &Workspace{
		Scripts: []Script{
			{
				ScruptID: "start",
				Title:    "Здравствуйте, я робот-ИФИ.",
				StartQID: "как зовут?",
				Options: []OptionOfScript{
					{
						Key:     "start",
						NextSID: "start",
					},
				},
			},
		},
		Questions: []Question{
			{
				QID: "как зовут?",
				Texts: []string{
					"a",
					"b",
				},
				ParamName: "user_name",
				NextQID:   "сколько лет?",
			},
			{
				QID:       "сколько лет?",
				Texts:     []string{"aaaa"},
				ParamName: "user_age",
				NextQID:   "человек это компьютер?",
				Options: []OptionOfResponse{
					{
						Key:     "yes",
						NextQID: "да",
					},
					{
						Key:     "no",
						NextQID: "нет",
					},
				},
			},
		},
	}
	w := &Workspace{}
	_, err := toml.Decode(dat, w)
	assert.NoError(t, err, "decode from toml")

	assert.Len(t, w.Scripts, 1)
	assert.Len(t, w.Questions, 2)

	assert.EqualValues(t, want, w)
}

var dat = `[[s]]
id = "start"
title = "Здравствуйте, я робот-ИФИ."
start_qid = "как зовут?"
	[[s.opts]]
	key = "start"
	sid = "start"

[[q]]
id = "как зовут?"
texts = [
    "a",
    "b",
]
param_name = "user_name"
next_qid = "сколько лет?"

[[q]]
id = "сколько лет?"
texts = [
    "aaaa",
]
param_name = "user_age"
next_qid = "человек это компьютер?"
[[q.opts]]
	key = "yes"
	qid = "да"
[[q.opts]]
	key = "no"
	qid = "нет"
`
