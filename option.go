package flatten

import "slices"

// Options configures Flatten.
type Options func(*option)

// WithIncludeFields includes specific fields in result.
func WithIncludeFields(fields ...string) Options {
	return func(opt *option) {
		for _, field := range fields {
			if !slices.Contains(opt.includes, field) {
				opt.includes = append(opt.includes, field)
			}
		}
	}
}

// WithExcludeFields excludes specific fields from result.
func WithExcludeFields(fields ...string) Options {
	return func(opt *option) {
		for _, field := range fields {
			if !slices.Contains(opt.excludes, field) {
				opt.excludes = append(opt.excludes, field)
			}
		}
	}
}

type option struct {
	includes []string
	excludes []string
}

func newOption() *option {
	return &option{
		includes: []string{},
		excludes: []string{},
	}
}

// shouldSkip determines whether a given field should be skipped.
func (opt *option) shouldSkip(field string) bool {
	if field != "" {
		if (len(opt.includes) > 0 && !slices.Contains(opt.includes, field)) ||
			(len(opt.excludes) > 0 && slices.Contains(opt.excludes, field)) {
			return true
		}
	}

	return false
}
