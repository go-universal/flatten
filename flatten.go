package flatten

import (
	"sort"
	"strings"
)

// Flatten converts a nested data structure into a flat slice of strings.
//
// Parameters:
//   - value: any nested data structure (map, slice, struct, or primitive)
//   - options: variadic Options functions to customize flattening behavior
//
// Returns:
//
//	A sorted slice of strings representing the flattened structure.
func Flatten(value any, options ...Options) []string {
	option := newOption()
	for _, o := range options {
		o(option)
	}

	var result []string
	collect(value, "", &result, false, *option)
	sort.Strings(result)
	return result
}

// FlattenCompare compares two values by flattening them and checking for equality.
// Returns true if the resulting strings are equal, returns false if the flattened representations differ.
func FlattenCompare(src, dest any, options ...Options) bool {
	source := strings.Join(Flatten(src, options...), "|")
	destination := strings.Join(Flatten(dest, options...), "|")
	return source == destination
}
