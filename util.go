package flatten

import (
	"fmt"
	"reflect"
	"strconv"
)

// isNil checks whether v is nil using reflect safely.
func isNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return true
	}

	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Interface:
		return rv.IsNil()
	}

	return false
}

// encode turns primitive values into string. Safe with reflect.
func encode(v any) string {
	if isNil(v) {
		return "[null]"
	}

	switch rv := reflect.Indirect(reflect.ValueOf(v)); rv.Kind() {
	case reflect.Invalid:
		return "[undefined]"

	case reflect.String:
		return rv.String()

	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 64)

	default:
		return fmt.Sprintf("%v", rv.Interface())
	}
}

// path generate dot separated path
func path(root, key string) string {
	if root != "" && key != "" {
		return root + "." + key
	} else if root != "" {
		return root
	} else if key != "" {
		return key
	}

	return ""
}

// flat make key:value flat string pair
func flat(key string, value any, isArray bool) string {
	if isArray {
		return fmt.Sprintf("%s:[%s]", key, encode(value))
	}

	return fmt.Sprintf("%s:%s", key, encode(value))
}
