# Go Flatten Library

![GitHub Tag](https://img.shields.io/github/v/tag/go-universal/flatten?sort=semver&label=version) [![Go Reference](https://pkg.go.dev/badge/github.com/go-universal/flatten.svg)](https://pkg.go.dev/github.com/go-universal/flatten) [![License](https://img.shields.io/badge/license-ISC-blue.svg)](https://github.com/go-universal/flatten/blob/main/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/go-universal/flatten)](https://goreportcard.com/report/github.com/go-universal/flatten) ![Contributors](https://img.shields.io/github/contributors/go-universal/flatten) ![Issues](https://img.shields.io/github/issues/go-universal/flatten)

A Go library for creating a normalized and flattened representation of any value. Perfect for data comparison, serialization, and deep structure analysis.

## Features

- **Flatten Any Type**: Support for primitives, arrays, maps, structs, and deeply nested structures
- **Field Filtering**: Include or exclude specific fields using options
- **Deep Comparison**: Compare complex structures regardless of order
- **Type Safe**: Works with Go's type system using reflection
- **Sorted Output**: Results are consistently sorted for reliable comparisons

## Installation

```sh
go get github.com/go-universal/flatten
```

## Quick Start

### Basic Flattening

```go
package main

import (
    "fmt"
    "github.com/go-universal/flatten"
)

type User struct {
    ID    int
    Name  string
    Email string
    age   int // unexported fields are ignored
}

func main() {
    user := User{
        ID:    1,
        Name:  "Alice",
        Email: "alice@example.com",
        age:   30,
    }

    result := flatten.Flatten(user)
    for _, item := range result {
        fmt.Println(item)
    }
    // Output:
    // Email:alice@example.com
    // ID:1
    // Name:Alice
}
```

### Filtering Fields

```go
// Include only specific fields
result := flatten.Flatten(user, flatten.WithIncludeFields("ID", "Name"))
// Result: [ID:1 Name:Alice]

// Exclude specific fields
result := flatten.Flatten(user, flatten.WithExcludeFields("Email"))
// Result: [ID:1 Name:Alice]
```

### Comparing Structures

```go
user1 := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
user2 := User{ID: 1, Name: "Alice", Email: "alice@example.com"}

// Compare all fields
isEqual := flatten.FlattenCompare(user1, user2)
fmt.Println(isEqual) // true

// Compare with excluded fields
user3 := User{ID: 1, Name: "Alice", Email: "different@example.com"}
isEqual = flatten.FlattenCompare(user1, user3, flatten.WithExcludeFields("Email"))
fmt.Println(isEqual) // true
```

## Supported Types

- **Primitives**: string, int, uint, float64, bool, nil
- **Collections**: arrays, slices, maps
- **Structures**: structs with exported fields
- **Nested**: arbitrary depth of nested types

## How It Works

The library recursively traverses your data structure and generates a flat list of key-value pairs:

```go
type Person struct {
    Name string
    Address struct {
        City string
        Zip  string
    }
}

person := Person{
    Name: "Bob",
    Address: struct {
        City string
        Zip  string
    }{
        City: "NYC",
        Zip:  "10001",
    },
}

result := flatten.Flatten(person)
// Output:
// Address.City:NYC
// Address.Zip:10001
// Name:Bob
```

## API Reference

### Flatten

```go
func Flatten(value any, options ...Options) []string
```

Flattens any value into a sorted slice of strings. Each string represents a flattened key-value pair in the format `key:value`.

**Parameters:**

- `value`: The value to flatten (any type)
- `options`: Optional configuration functions

**Returns:** Sorted slice of flattened strings

### FlattenCompare

```go
func FlattenCompare(src, dest any, options ...Options) bool
```

Compares two values by flattening them and checking for equality.

**Parameters:**

- `src`: Source value to compare
- `dest`: Destination value to compare
- `options`: Optional configuration functions

**Returns:** `true` if both values flatten to the same representation

### Options

#### WithIncludeFields

```go
flatten.WithIncludeFields(fields ...string) Options
```

Includes only specified fields in the flattening result.

#### WithExcludeFields

```go
flatten.WithExcludeFields(fields ...string) Options
```

Excludes specified fields from the flattening result.

## Examples

### Example 1: Flatten Complex Nested Structure

```go
type Company struct {
    Name string
    Employees []struct {
        ID   int
        Name string
    }
}

company := Company{
    Name: "TechCorp",
    Employees: []struct {
        ID   int
        Name string
    }{
        {ID: 1, Name: "Alice"},
        {ID: 2, Name: "Bob"},
    },
}

result := flatten.Flatten(company)
// Output:
// Employees.ID:1
// Employees.ID:2
// Employees.Name:Alice
// Employees.Name:Bob
// Name:TechCorp
```

### Example 2: Compare Structures Ignoring Timestamps

```go
type Document struct {
    Title      string
    Content    string
    ModifiedAt time.Time
}

doc1 := Document{Title: "Report", Content: "Data", ModifiedAt: time.Now()}
doc2 := Document{Title: "Report", Content: "Data", ModifiedAt: time.Now().Add(1 * time.Hour)}

// Compare documents ignoring the timestamp
isEqual := flatten.FlattenCompare(doc1, doc2, flatten.WithExcludeFields("ModifiedAt"))
fmt.Println(isEqual) // true
```

## License

This project is licensed under the ISC License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
