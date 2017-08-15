package utils

import (
	"bytes"
	"encoding/json"
	"log"
)

type ValueType string

func (v ValueType) String() string {
	return string(v)
}

var (
	TUnknown   ValueType = ""
	TString    ValueType = "string"
	TStringArr ValueType = "string_arr"
	TInteger   ValueType = "integer"
	TFloat     ValueType = "float"
	TBool      ValueType = "bool"
)

func NewValues() *Values {
	return &Values{
		S:    make(map[string]string),
		SArr: make(map[string][]string),
		I:    make(map[string]int64),
		F:    make(map[string]float64),
		B:    make(map[string]bool),
	}
}

type Values struct {
	S    map[string]string
	SArr map[string][]string
	I    map[string]int64
	F    map[string]float64
	B    map[string]bool
}

type value struct {
	Value interface{} `json:"value"`
	Type  ValueType   `json:"type"`
}

func (v Values) MarshalJSON() (dat []byte, err error) {
	res := map[string]value{}
	for k, v := range v.S {
		res[k] = value{v, TString}
	}
	for k, v := range v.SArr {
		res[k] = value{v, TStringArr}
	}
	for k, v := range v.I {
		res[k] = value{v, TInteger}
	}
	for k, v := range v.F {
		res[k] = value{v, TFloat}
	}
	for k, v := range v.B {
		res[k] = value{v, TBool}
	}
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(res)
	return buf.Bytes(), err
}

func (v *Values) UnmarshalJSON(dat []byte) (err error) {
	skelet := map[string]value{}
	if err := json.NewDecoder(bytes.NewBuffer(dat)).Decode(&skelet); err != nil {
		return err
	}

	for k, _v := range skelet {

		// TODO: приведение типов более тщательное

		switch _v.Type {
		case TString:
			v.S[k] = _v.Value.(string)
		case TStringArr:
			vt, ok := _v.Value.([]interface{})
			if !ok {
				log.Printf("expected []interface{}, got %T\n", _v.Value)
				continue
			}
			arr := make([]string, len(vt))
			for _i, _vv := range vt {
				arr[_i] = _vv.(string)
			}
			v.SArr[k] = arr
		case TInteger:
			v.I[k] = int64(_v.Value.(float64))
		case TFloat:
			v.F[k] = _v.Value.(float64)
		case TBool:
			v.B[k] = _v.Value.(bool)
		default:
			log.Printf("not supported value type %q (%q: %+v)\n", _v.Type, k, _v.Value)
		}
	}

	return
}
