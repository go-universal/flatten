package flatten

import (
	"reflect"
)

type Transformer func(v any) []string

var transformers = make(map[reflect.Type]Transformer)

func RegisterTransformer[T any](fn func(v T) []string) {
	var zero T
	t := reflect.TypeOf(zero)

	if t == nil {
		return
	}

	transformers[t] = func(v any) []string {
		return fn(v.(T))
	}
}

func resolve(v any) ([]string, bool) {
	if v == nil {
		return nil, false
	}

	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)

	// pointer-safe
	if rt.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, false
		}
		rt = rt.Elem()
		v = rv.Elem().Interface()
	}

	if fn, ok := transformers[rt]; ok {
		return fn(v), true
	}

	return nil, false
}
