package utils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestField2Marshal(t *testing.T) {
	v := Values{
		S: map[string]string{
			"str1": "foo",
			"str2": "bar",
		},
		SArr: map[string][]string{
			"arr1": []string{"tag1", "tag2"},
		},
		I: map[string]int64{
			"int1": 123,
			"int2": 456,
		},
		F: map[string]float64{
			"f1": 123.456,
			"f2": 456.789,
		},
		B: map[string]bool{
			"b1": true,
			"b2": false,
		},
	}

	gotJSON, err := json.Marshal(v)

	expectedJSON := `{"arr1":{"value":["tag1","tag2"],"type":"string_arr"},"b1":{"value":true,"type":"bool"},"b2":{"value":false,"type":"bool"},"f1":{"value":123.456,"type":"float"},"f2":{"value":456.789,"type":"float"},"int1":{"value":123,"type":"integer"},"int2":{"value":456,"type":"integer"},"str1":{"value":"foo","type":"string"},"str2":{"value":"bar","type":"string"}}`

	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(gotJSON))
}

func TestField2Unmarshal(t *testing.T) {
	dat := `{"arr1":{"value":["tag1","tag2"],"type":"string_arr"},"b1":{"value":true,"type":"bool"},"b2":{"value":false,"type":"bool"},"f1":{"value":123.456,"type":"float"},"f2":{"value":456.789,"type":"float"},"int1":{"value":123,"type":"integer"},"int2":{"value":456,"type":"integer"},"str1":{"value":"foo","type":"string"},"str2":{"value":"bar","type":"string"}}`
	got := Values{
		S: map[string]string{
			"str1": "foo",
			"str2": "bar",
		},
		SArr: map[string][]string{
			"arr1": []string{"tag1", "tag2"},
		},
		I: map[string]int64{
			"int1": 123,
			"int2": 456,
		},
		F: map[string]float64{
			"f1": 123.456,
			"f2": 456.789,
		},
		B: map[string]bool{
			"b1": true,
			"b2": false,
		},
	}

	v := NewValues()
	err := json.Unmarshal([]byte(dat), &v)
	assert.NoError(t, err)
	assert.Equal(t, got.S, v.S)
	assert.Equal(t, got.SArr, v.SArr)
	assert.Equal(t, got.I, v.I)
	assert.Equal(t, got.F, v.F)
	assert.Equal(t, got.B, v.B)
}
