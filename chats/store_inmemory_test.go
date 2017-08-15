package chats

import "testing"
import "github.com/stretchr/testify/assert"

func TestInmemoryState_Update(t *testing.T) {
	type args struct {
		obj *State
	}
	tests := []struct {
		name      string
		fields    *InmemoryState
		args      args
		wantErr   bool
		wantState *State
	}{
		{
			"",
			&InmemoryState{
				List: map[string]State{
					"1": State{
						ChatID:   "1",
						ScriptID: "start",
						LastQID:  "q#1",
						Props: map[string]interface{}{
							"a": "b",
							"c": "d",
						},
					},
				},
			},
			args{
				&State{
					ChatID:   "1",
					ScriptID: "start",
					LastQID:  "q#2",
					Props: map[string]interface{}{
						"a": "b",
						"e": "f",
					},
				},
			},
			false,
			&State{
				ChatID:   "1",
				ScriptID: "start",
				LastQID:  "q#2",
				Props: map[string]interface{}{
					"a": "b",
					"e": "f",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &InmemoryState{
				List: tt.fields.List,
			}
			if err := s.Update(tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("InmemoryState.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, _ := s.Find(tt.args.obj.ChatID)
			assert.Equal(t, tt.wantState, got)
		})
	}
}
