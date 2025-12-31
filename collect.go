package flatten

import (
	"fmt"
	"reflect"
)

// collect recursively traverses the input value and flattens its structure into a slice of strings.
func collect(inp any, prefix string, out *[]string, isArray bool, opt option) {
	// Skip
	if opt.shouldSkip(prefix) {
		return
	}

	// Skip Nil Fields
	if isNil(inp) {
		*out = append(*out, flat(prefix, nil, isArray))
		return
	}

	// Transformers
	if res, ok := resolve(inp); ok {
		for _, item := range res {
			*out = append(*out, flat(prefix, item, isArray))
		}

		return
	}

	value := reflect.Indirect(reflect.ValueOf(inp))
	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			collect(value.Index(i).Interface(), prefix, out, true, opt)
		}
	case reflect.Map:
		for _, k := range value.MapKeys() {
			fullPath := path(prefix, fmt.Sprint(k.Interface()))
			collect(value.MapIndex(k).Interface(), fullPath, out, false, opt)
		}
	case reflect.Struct:
		typ := value.Type()
		for i := 0; i < value.NumField(); i++ {
			field := typ.Field(i)
			name := field.Name
			fullPath := path(prefix, name)

			if !field.IsExported() {
				continue
			}

			collect(value.Field(i).Interface(), fullPath, out, false, opt)
		}
	default:
		*out = append(*out, flat(prefix, inp, isArray))
	}
}
