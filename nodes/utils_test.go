package nodes

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func TestTexts_UnmarshalTOML(t *testing.T) {
	type dto struct {
		Text Texts `toml:"text"`
	}
	tests := []struct {
		data    string
		want    Texts
		wantErr bool
	}{
		{
			`text = ["a", "b", "c"]`,
			Texts{"a", "b", "c"},
			false,
		},
		{
			`text = "a"`,
			Texts{"a"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := dto{}
			err := toml.Unmarshal([]byte(tt.data), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Texts.UnmarshalTOML() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.EqualValues(t, tt.want, got.Text)
		})
	}
}
