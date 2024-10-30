package db

import (
	"reflect"
)

type NamedEntity interface {
	GetValue() string
}

func ValuesFromAny(entities interface{}) []string {
	v := reflect.ValueOf(entities)

	if v.Kind() != reflect.Slice {
		panic("expected a slice")
	}

	values := make([]string, v.Len())

	for i := 0; i < v.Len(); i++ {
		entity := v.Index(i).Interface().(NamedEntity)
		values[i] = entity.GetValue()
	}

	return values
}
