package utils

import (
	"encoding/binary"
	"log"
	"math"

	"reflect"

	lua "github.com/yuin/gopher-lua"
)

// ToValueFromLValue преобразует lua.LValue в основной тип данных
// Возвращаемые типы:
// - LNumber -> float64
// - LBool -> bool
// - LString -> string
// - LTable -> []interface{} or map[string]interface{}
// - LNil -> nil
func ToValueFromLValue(v lua.LValue) interface{} {
	switch v.Type() {
	case lua.LTNumber:
		return float64(v.(lua.LNumber))
	case lua.LTBool:
		return bool(v.(lua.LBool))
	case lua.LTString:
		return string(v.(lua.LString))
	case lua.LTUserData:
		return v.(*lua.LUserData).Value
	case lua.LTNil:
		return nil
	case lua.LTTable:
		tbl := v.(*lua.LTable)
		var keys []string
		var vals []interface{}

		isArray := true
		counter := 0
		tbl.ForEach(func(k, v lua.LValue) {
			if k.Type() == lua.LTNumber && int(k.(lua.LNumber)) != counter+1 {
				isArray = false
			}
			if k.Type() != lua.LTNumber {
				isArray = false
			}

			keys = append(keys, k.String())
			vals = append(vals, ToValueFromLValue(v))

			counter++
		})

		if isArray {
			return vals
		}

		_vals := make(map[string]interface{}, counter)
		for i := 0; i < counter; i++ {
			_vals[keys[i]] = vals[i]
		}

		return _vals
	default:
		log.Println("not supported type", v.Type())
	}

	return nil
}

func ToLValueOrNil(v interface{}, L *lua.LState) lua.LValue {
	switch v := v.(type) {
	case reflect.Value:
		if v.IsValid() {
			return ToLValueOrNil(v.Interface(), L)
		}
	case int, int8, int32, int64, float32, float64,
		uint, uint8, uint16, uint32, uint64:

		return lua.LNumber(ToFloat64(v))
	case bool:
		return lua.LBool(v)
	case string:
		return lua.LString(v)
	case nil:
		return lua.LNil
	case []int, []int64, []float64, []string, []bool, []interface{}:
		tb := L.NewTable()

		var litems []lua.LValue

		// types
		switch v := v.(type) {
		case []int:
			litems = make([]lua.LValue, len(v))
			for index, item := range v {
				if _v := ToLValueOrNil(item, L); _v != nil {
					litems[index] = _v
				}
			}
		case []int64:
			litems = make([]lua.LValue, len(v))
			for index, item := range v {
				if _v := ToLValueOrNil(item, L); _v != nil {
					litems[index] = _v
				}
			}
		case []float64:
			litems = make([]lua.LValue, len(v))
			for index, item := range v {
				if _v := ToLValueOrNil(item, L); _v != nil {
					litems[index] = _v
				}
			}
		case []bool:
			litems = make([]lua.LValue, len(v))
			for index, item := range v {
				if _v := ToLValueOrNil(item, L); _v != nil {
					litems[index] = _v
				}
			}
		case []string:
			litems = make([]lua.LValue, len(v))
			for index, item := range v {
				if _v := ToLValueOrNil(item, L); _v != nil {
					litems[index] = _v
				}
			}
		case []interface{}:
			litems = make([]lua.LValue, len(v))
			for index, item := range v {
				if _v := ToLValueOrNil(item, L); _v != nil {
					litems[index] = _v
				}
			}
		default:
			log.Printf(
				"[ERR] ToLValueOrNil: slice, not expected type value, got %T",
				v,
			)
		}

		// main

		if len(litems) != 0 {
			for _, item := range litems {
				tb.Append(item)
			}
		}

		return tb

	case map[string]interface{}, map[interface{}]interface{}:
		tb := L.NewTable()

		var keys, values []lua.LValue

		// types
		switch v := v.(type) {
		case map[string]string:
		case map[string]interface{}:
			keys = make([]lua.LValue, len(v))
			values = make([]lua.LValue, len(v))
			var seq = 0
			for key, value := range v {
				keys[seq] = ToLValueOrNil(key, L)
				values[seq] = ToLValueOrNil(value, L)
				seq++
			}
		case map[interface{}]interface{}:
			keys = make([]lua.LValue, len(v))
			values = make([]lua.LValue, len(v))
			var seq = 0
			for key, value := range v {
				keys[seq] = ToLValueOrNil(key, L)
				values[seq] = ToLValueOrNil(value, L)
				seq++
			}
		default:
			log.Printf(
				"[ERR] ToLValueOrNil: not expected type value, got %T",
				v,
			)
		}

		for i := 0; i < len(keys); i++ {
			tb.RawSet(keys[i], values[i])
		}

		return tb
	default:
		log.Printf(
			"[ERR] ToLValueOrNil: not expected type value, got %T",
			v,
		)
	}

	return lua.LNil
}

func ToFloat64(v interface{}) (f float64) {
	switch _v := v.(type) {
	case int:
		f = float64(_v)
	case int16:
		f = float64(_v)
	case int32:
		f = float64(_v)
	case int64:
		f = float64(_v)
	case int8:
		f = float64(_v)
	case float32:
		f = float64(_v)
	case float64:
		f = float64(_v)
	case uint:
		f = float64(_v)
	case uint16:
		f = float64(_v)
	case uint32:
		f = float64(_v)
	case uint64:
		f = float64(_v)
	case uint8:
		f = float64(_v)
	default:
		f = 0
	}

	return
}

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}
