package nodes

import (
	"reflect"
	"strconv"
)

type Texts []string

func (t *Texts) UnmarshalTOML(in interface{}) error {
	vt := reflect.TypeOf(in)
	v := reflect.ValueOf(in)
	switch vt.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			riv := v.Index(i)
			if riv.Type().Kind() != reflect.Interface {
				continue
			}

			var str string
			switch _v := riv.Interface().(type) {
			case string:
				str = _v
			case int64:
				str = strconv.FormatInt(_v, 64)
			case float64:
				str = strconv.FormatFloat(_v, 'E', -1, 64)
			default:
				// not supported type
				continue
			}

			*t = append(*t, str)
		}
	case reflect.String:
		*t = Texts{v.Interface().(string)}
	}
	return nil
}
