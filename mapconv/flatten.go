package mapconv

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

func FlattenRec(v reflect.Value, prefix string, omit bool, m map[string]interface{}) error {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v = reflect.Zero(v.Type().Elem())
			break
		}
		v = v.Elem()
	}
	t := v.Type()
fieldLoop:
	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)
		name := ft.Name
		if !unicode.IsUpper(rune(name[0])) {
			continue
		}
		if tag, ok := ft.Tag.Lookup("json"); ok {
			name = strings.Split(tag, ",")[0]
			if name == "-" {
				continue
			}
		}
		fv := v.Field(i)
		for fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				if omit {
					continue fieldLoop
				}
				fv = reflect.Zero(fv.Type().Elem())
				break
			}
			fv = fv.Elem()
		}
		if fvt := fv.Type(); fvt == TimeType || fvt == DecimalType || fvt == UUIDType {
			m[name] = fmt.Sprint(fv.Interface())
		} else if fv.Kind() == reflect.Struct {
			FlattenRec(fv, prefix+name+".", omit, m)
		} else {
			m[prefix+name] = fv.Interface()
		}
	}
	return nil
}

func Flatten(x interface{}, omit bool) (map[string]interface{}, error) {
	t := reflect.TypeOf(x)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input type should be struct, got %T", x)
	}
	m := map[string]interface{}{}
	err := FlattenRec(reflect.ValueOf(x), "", omit, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
