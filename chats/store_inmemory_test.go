package chats

import "testing"
import "github.com/stretchr/testify/assert"

func TestInmemoryChat_Update(t *testing.T) {
	type args struct {
		obj *Chat
	}
	tests := []struct {
		name     string
		fields   *InmemoryChat
		args     args
		wantErr  bool
		wantChat *Chat
	}{
		{
			"",
			&InmemoryChat{
				List: map[string]Chat{
					"1": Chat{
						ChatID:         "1",
						PreviousNodeID: "start",
						NextNodeID:     "q#1",
						Props: map[string]interface{}{
							"a": "b",
							"c": "d",
						},
					},
				},
			},
			args{
				&Chat{
					ChatID:         "1",
					PreviousNodeID: "start",
					NextNodeID:     "q#2",
					Props: map[string]interface{}{
						"a": "b",
						"e": "f",
					},
				},
			},
			false,
			&Chat{
				ChatID:         "1",
				PreviousNodeID: "start",
				NextNodeID:     "q#2",
				Props: map[string]interface{}{
					"a": "b",
					"e": "f",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &InmemoryChat{
				List: tt.fields.List,
			}
			if err := s.Update(tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("InmemoryChat.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, _ := s.Find(tt.args.obj.ChatID)
			assert.Equal(t, tt.wantChat.ChatID, got.ChatID)
			assert.Equal(t, tt.wantChat.PreviousNodeID, got.PreviousNodeID)
			assert.Equal(t, tt.wantChat.NextNodeID, got.NextNodeID)
			assert.Equal(t, tt.wantChat.Props, got.Props)
		})
	}
}
