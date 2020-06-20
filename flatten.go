package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/shopspring/decimal"
)

var timeType = reflect.TypeOf(date.Time{})
var decimalType = reflect.TypeOf(decimal.Decimal{})

func flattenToMapRec(v reflect.Value, omit bool, m map[string]interface{}) error {
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
		if tag, ok := ft.Tag.Lookup("json"); ok {
			name = strings.Split(tag, ",")[0]
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
		if fv.Type() != timeType && fv.Type() != decimalType && fv.Kind() == reflect.Struct {
			flattenToMapRec(fv, omit, m)
		} else {
			m[name] = fv.Interface()
		}
	}
	return nil
}

func flattenToMap(x interface{}, omit bool) (map[string]interface{}, error) {
	t := reflect.TypeOf(x)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input type should be struct, got %T", x)
	}
	m := map[string]interface{}{}
	err := flattenToMapRec(reflect.ValueOf(x), omit, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
