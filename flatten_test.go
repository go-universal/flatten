package flatten

import (
	"fmt"
	"testing"
)

func TestResolver(t *testing.T) {
	type CustomType struct {
		name string
		Age  int
	}

	tests := []struct {
		name     string
		input    any
		expected []string
	}{
		{
			name:     "resolver",
			input:    CustomType{name: "Alice", Age: 30},
			expected: []string{"name:Alice"},
		},
	}

	RegisterTransformer(func(v CustomType) []string {
		return []string{
			fmt.Sprintf("name:%s", v.name),
		}
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Flatten(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("length mismatch:\n    got (%d): %+v \n    expected (%d) %+v", len(result), result, len(tt.expected), tt.expected)
				return
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected []string
	}{
		{
			name:     "primitive string",
			input:    "hello",
			expected: []string{":hello"},
		},
		{
			name:     "primitive int",
			input:    42,
			expected: []string{":42"},
		},
		{
			name:     "primitive bool",
			input:    true,
			expected: []string{":true"},
		},
		{
			name:     "primitive float",
			input:    3.14,
			expected: []string{":3.14"},
		},
		{
			name:     "nil value",
			input:    nil,
			expected: []string{":[null]"},
		},
		{
			name:     "simple array",
			input:    []int{1, 2, 3},
			expected: []string{":[1]", ":[2]", ":[3]"},
		},
		{
			name:     "simple map",
			input:    map[string]int{"a": 1, "b": 2},
			expected: []string{"a:1", "b:2"},
		},
		{
			name: "simple struct",
			input: struct {
				Name string
				Age  int
			}{
				Name: "Alice",
				Age:  30,
			},
			expected: []string{"Age:30", "Name:Alice"},
		},
		{
			name: "struct with array",
			input: struct {
				Name    string
				Hobbies []string
			}{
				Name:    "Bob",
				Hobbies: []string{"reading", "gaming"},
			},
			expected: []string{"Hobbies:[gaming]", "Hobbies:[reading]", "Name:Bob"},
		},
		{
			name: "struct with map",
			input: struct {
				Name   string
				Scores map[string]int
			}{
				Name:   "Charlie",
				Scores: map[string]int{"math": 95, "english": 87},
			},
			expected: []string{"Name:Charlie", "Scores.english:87", "Scores.math:95"},
		},
		{
			name: "deeply nested structure",
			input: struct {
				User struct {
					Profile struct {
						Name    string
						Contact struct {
							Email string
							Phone string
						}
					}
				}
			}{
				User: struct {
					Profile struct {
						Name    string
						Contact struct {
							Email string
							Phone string
						}
					}
				}{
					Profile: struct {
						Name    string
						Contact struct {
							Email string
							Phone string
						}
					}{
						Name: "David",
						Contact: struct {
							Email string
							Phone string
						}{
							Email: "david@example.com",
							Phone: "555-1234",
						},
					},
				},
			},
			expected: []string{
				"User.Profile.Contact.Email:david@example.com",
				"User.Profile.Contact.Phone:555-1234",
				"User.Profile.Name:David",
			},
		},
		{
			name: "array of structs",
			input: []struct {
				ID   int
				Name string
			}{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
			expected: []string{"ID:1", "ID:2", "Name:Alice", "Name:Bob"},
		},
		{
			name: "complex nested with arrays and maps",
			input: struct {
				Items []struct {
					Name string
					Tags map[string]string
				}
			}{
				Items: []struct {
					Name string
					Tags map[string]string
				}{
					{
						Name: "Item1",
						Tags: map[string]string{"color": "red", "size": "large"},
					},
					{
						Name: "Item2",
						Tags: map[string]string{"color": "blue"},
					},
				},
			},
			expected: []string{
				"Items.Name:Item1",
				"Items.Name:Item2",
				"Items.Tags.color:blue",
				"Items.Tags.color:red",
				"Items.Tags.size:large",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Flatten(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("length mismatch: got %d, expected %d\nGot: %v\nExpected: %v", len(result), len(tt.expected), result, tt.expected)
				return
			}

			for i, exp := range tt.expected {
				if result[i] != exp {
					t.Errorf("mismatch at index %d: got %q, expected %q", i, result[i], exp)
				}
			}
		})
	}
}

func TestFlattenWithOptions(t *testing.T) {
	user := struct {
		ID    int
		Name  string
		Email string
	}{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
	}

	t.Run("with include fields", func(t *testing.T) {
		result := Flatten(user, WithIncludeFields("Name"))
		expected := []string{"Name:Alice"}

		if len(result) != len(expected) {
			t.Errorf("length mismatch: got %d, expected %d", len(result), len(expected))
			return
		}

		if result[0] != expected[0] {
			t.Errorf("got %q, expected %q", result[0], expected[0])
		}
	})

	t.Run("with exclude fields", func(t *testing.T) {
		result := Flatten(user, WithExcludeFields("Email"))
		expected := []string{"ID:1", "Name:Alice"}

		if len(result) != len(expected) {
			t.Errorf("length mismatch: got %d, expected %d", len(result), len(expected))
			return
		}

		for i, exp := range expected {
			if result[i] != exp {
				t.Errorf("got %q, expected %q", result[i], exp)
			}
		}
	})
}

func TestFlattenCompare(t *testing.T) {
	tests := []struct {
		name     string
		src      any
		dest     any
		expected bool
	}{
		{
			name:     "identical primitives",
			src:      42,
			dest:     42,
			expected: true,
		},
		{
			name:     "different primitives",
			src:      42,
			dest:     43,
			expected: false,
		},
		{
			name:     "identical strings",
			src:      "hello",
			dest:     "hello",
			expected: true,
		},
		{
			name:     "different strings",
			src:      "hello",
			dest:     "world",
			expected: false,
		},
		{
			name:     "identical arrays",
			src:      []int{1, 2, 3},
			dest:     []int{1, 2, 3},
			expected: true,
		},
		{
			name:     "different array lengths",
			src:      []int{1, 2, 3},
			dest:     []int{1, 2},
			expected: false,
		},
		{
			name:     "different array values",
			src:      []int{1, 2, 3},
			dest:     []int{1, 2, 4},
			expected: false,
		},
		{
			name:     "identical maps",
			src:      map[string]int{"a": 1, "b": 2},
			dest:     map[string]int{"a": 1, "b": 2},
			expected: true,
		},
		{
			name:     "different map values",
			src:      map[string]int{"a": 1, "b": 2},
			dest:     map[string]int{"a": 1, "b": 3},
			expected: false,
		},
		{
			name:     "different map keys",
			src:      map[string]int{"a": 1, "b": 2},
			dest:     map[string]int{"a": 1, "c": 2},
			expected: false,
		},
		{
			name: "identical structs",
			src: struct {
				Name string
				Age  int
			}{
				Name: "Alice",
				Age:  30,
			},
			dest: struct {
				Name string
				Age  int
			}{
				Name: "Alice",
				Age:  30,
			},
			expected: true,
		},
		{
			name: "different struct values",
			src: struct {
				Name string
				Age  int
			}{
				Name: "Alice",
				Age:  30,
			},
			dest: struct {
				Name string
				Age  int
			}{
				Name: "Bob",
				Age:  30,
			},
			expected: false,
		},
		{
			name: "identical nested structs",
			src: struct {
				User struct {
					Name  string
					Email string
				}
			}{
				User: struct {
					Name  string
					Email string
				}{
					Name:  "Charlie",
					Email: "charlie@example.com",
				},
			},
			dest: struct {
				User struct {
					Name  string
					Email string
				}
			}{
				User: struct {
					Name  string
					Email string
				}{
					Name:  "Charlie",
					Email: "charlie@example.com",
				},
			},
			expected: true,
		},
		{
			name: "different nested struct values",
			src: struct {
				User struct {
					Name  string
					Email string
				}
			}{
				User: struct {
					Name  string
					Email string
				}{
					Name:  "Charlie",
					Email: "charlie@example.com",
				},
			},
			dest: struct {
				User struct {
					Name  string
					Email string
				}
			}{
				User: struct {
					Name  string
					Email string
				}{
					Name:  "Charlie",
					Email: "charlie@example.org",
				},
			},
			expected: false,
		},
		{
			name: "identical complex structures",
			src: struct {
				Items []struct {
					ID   int
					Tags map[string]string
				}
			}{
				Items: []struct {
					ID   int
					Tags map[string]string
				}{
					{
						ID:   1,
						Tags: map[string]string{"color": "red", "size": "large"},
					},
				},
			},
			dest: struct {
				Items []struct {
					ID   int
					Tags map[string]string
				}
			}{
				Items: []struct {
					ID   int
					Tags map[string]string
				}{
					{
						ID:   1,
						Tags: map[string]string{"color": "red", "size": "large"},
					},
				},
			},
			expected: true,
		},
		{
			name: "different complex structures",
			src: struct {
				Items []struct {
					ID   int
					Tags map[string]string
				}
			}{
				Items: []struct {
					ID   int
					Tags map[string]string
				}{
					{
						ID:   1,
						Tags: map[string]string{"color": "red", "size": "large"},
					},
				},
			},
			dest: struct {
				Items []struct {
					ID   int
					Tags map[string]string
				}
			}{
				Items: []struct {
					ID   int
					Tags map[string]string
				}{
					{
						ID:   2,
						Tags: map[string]string{"color": "red", "size": "large"},
					},
				},
			},
			expected: false,
		},
		{
			name:     "both nil values",
			src:      nil,
			dest:     nil,
			expected: true,
		},
		{
			name:     "nil vs non-nil",
			src:      nil,
			dest:     "value",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FlattenCompare(tt.src, tt.dest)
			if result != tt.expected {
				t.Errorf("FlattenCompare(%v, %v) = %v, expected %v", tt.src, tt.dest, result, tt.expected)
			}
		})
	}
}

func TestFlattenCompareWithOptions(t *testing.T) {
	user1 := struct {
		ID    int
		Name  string
		Email string
	}{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
	}

	user2 := struct {
		ID    int
		Name  string
		Email string
	}{
		ID:    1,
		Name:  "Alice",
		Email: "alice@different.com",
	}

	t.Run("compare with include fields - same included fields", func(t *testing.T) {
		result := FlattenCompare(user1, user2, WithIncludeFields("ID", "Name"))
		if !result {
			t.Errorf("expected true when comparing included fields that are the same")
		}
	})

	t.Run("compare with include fields - different included fields", func(t *testing.T) {
		user3 := struct {
			ID    int
			Name  string
			Email string
		}{
			ID:    1,
			Name:  "Bob",
			Email: "alice@example.com",
		}
		result := FlattenCompare(user1, user3, WithIncludeFields("ID", "Name"))
		if result {
			t.Errorf("expected false when comparing included fields that are different")
		}
	})

	t.Run("compare with exclude fields - same after excluding", func(t *testing.T) {
		result := FlattenCompare(user1, user2, WithExcludeFields("Email"))
		if !result {
			t.Errorf("expected true when comparing fields after excluding the different field")
		}
	})

	t.Run("compare with exclude fields - different after excluding", func(t *testing.T) {
		user3 := struct {
			ID    int
			Name  string
			Email string
		}{
			ID:    1,
			Name:  "Bob",
			Email: "alice@example.com",
		}
		result := FlattenCompare(user1, user3, WithExcludeFields("Email"))
		if result {
			t.Errorf("expected false when comparing fields after excluding, but other fields differ")
		}
	})
}
