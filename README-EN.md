# gogeneric

[简体中文](README.md) | English

[![Go Reference](https://pkg.go.dev/badge/github.com/liu-houliang/gogeneric/v2.svg)](https://pkg.go.dev/github.com/liu-houliang/gogeneric/v2)
[![CI](https://github.com/liu-houliang/gogeneric/actions/workflows/ci.yml/badge.svg)](https://github.com/liu-houliang/gogeneric/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)](https://go.dev/)

A Go generics utility library that provides functional programming primitives and business-layer helpers that the standard library (`slices`/`maps`) intentionally omits.

**Requires Go >= 1.21**

## Installation

```bash
go get -u github.com/liu-houliang/gogeneric/v2
```

## Features

### 🔧 Functional Slice Helpers (`generic.go`)

These functions are **not provided by the standard library** and are the core value of this package.

```go
// Filter: keep elements matching a predicate
evens := gogeneric.Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
// → [2 4]

// Map: transform each element
strs := gogeneric.Map([]int{1, 2, 3}, func(n int) string { return fmt.Sprintf("%d", n) })
// → ["1" "2" "3"]

// Reduce: fold a slice into a single value
sum := gogeneric.Reduce([]int{1, 2, 3, 4, 5}, 0, func(acc, n int) int { return acc + n })
// → 15

// ForEach: iterate and perform an action on each element (no return value)
gogeneric.ForEach([]string{"a", "b"}, func(s string) { fmt.Println(s) })

// Chunk: split a slice into fixed-size sub-slices
gogeneric.Chunk([]int{1, 2, 3, 4, 5}, 2)
// → [[1 2] [3 4] [5]]

// Flatten: merge a slice of slices into one
gogeneric.Flatten([][]int{{1, 2}, {3, 4}, {5}})
// → [1 2 3 4 5]

// Unique: remove duplicates, preserving first-occurrence order
gogeneric.Unique([]int{1, 2, 2, 3, 1})
// → [1 2 3]
```

### 🧩 Generic Set Type (`set.go`)

```go
s := gogeneric.NewSet(1, 2, 3)
s.Add(4)
s.Contains(4)          // true
s.Remove(1)

a := gogeneric.NewSet(1, 2, 3)
b := gogeneric.NewSet(2, 3, 4)
a.Union(b)             // {1,2,3,4}
a.Intersection(b)      // {2,3}
a.Difference(b)        // {1}
a.IsSubset(b)          // false
```

### 🎯 Optional Type (`optional.go`)

Go has no built-in Option/Maybe type. This provides a type-safe implementation:

```go
opt := gogeneric.Some("hello")
opt.IsPresent()              // true
opt.GetOrElse("world")       // "hello"

empty := gogeneric.None[string]()
empty.GetOrElse("world")     // "world"
_, err := empty.GetOrError("not found") // err != nil

// Chainable transforms
length := gogeneric.MapOptional(gogeneric.Some("hello"),
    func(s string) int { return len(s) })
length.GetOrElse(0)          // 5
```

### 🔧 Struct Utilities (`generic.go`)

```go
type Student struct { Name string; Age uint }
students := []Student{{"alice", 20}, {"bob", 25}}

// Convert a struct slice to a map keyed by a field
m := gogeneric.StructSliceToMap[string]("Name", students)
// → map["alice":{alice 20} "bob":{bob 25}]

// Extract a single field from each struct into a slice
names := gogeneric.GetFieldArray[string]("Name", students)
// → ["alice" "bob"]

// Generic ternary operator
result := gogeneric.Cond(score > 60, "pass", "fail")
```

### 🌐 HTTP Utilities (`util.go`)

All HTTP functions return `*HTTPResponse[T]`, giving you the decoded body, status code, and response headers:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// GET — T can be struct, slice, map, pointer, or any type json.Unmarshal supports
res, err := gogeneric.HttpGetJson[[]User](ctx, "https://api.example.com/users", nil)
res.Data        // []User — decoded body
res.StatusCode  // 200
res.Header.Get("X-RateLimit-Remaining")  // response headers

// POST JSON body
res, err := gogeneric.HttpPostJson[User](ctx, url, CreateReq{Name: "alice"}, nil)

// PUT / PATCH JSON body
res, err := gogeneric.HttpPutJson[User](ctx, url+"/1", UpdateReq{Name: "bob"}, nil)
res, err := gogeneric.HttpPatchJson[User](ctx, url+"/1", map[string]string{"name": "bob"}, nil)

// DELETE (supports 204 No Content — Data will be zero value)
res, err := gogeneric.HttpDeleteJson[DeletedResp](ctx, url+"/1", nil)

// POST form-encoded (OAuth, legacy APIs)
form := url.Values{"grant_type": {"password"}, "username": {"alice"}}
res, err := gogeneric.HttpPostForm[TokenResp](ctx, "https://auth.example.com/token", form, nil)

// File upload via multipart/form-data
files := []gogeneric.File{
    {FieldName: "avatar", FileName: "photo.jpg", Reader: f},
}
res, err := gogeneric.HttpPostMultipart[UploadResp](ctx, url, files, map[string]string{"userId": "1"}, nil)
```

### ⚠️ Deprecated Functions (kept for backward compatibility)

The following functions are now covered by the Go 1.21+ standard library.
Migration is recommended:

| This library | Standard library replacement |
|-------------|------------------------------|
| `Min` / `Max` | built-in `min()` / `max()` |
| `SortSlice` | `slices.SortFunc` |
| `SliceIndex` | `slices.Index` |
| `CompareSlice` | `slices.Equal` |
| `MapKeys` | `maps.Keys` |
| `MapToSlice` | `maps.Values` |

## Running Tests

```bash
git clone https://github.com/liu-houliang/gogeneric.git
cd gogeneric

go test -v .                      # all tests with verbose output
go test -run TestFilter .         # only Filter tests
go test -run TestHttp .           # only HTTP tests
go test -bench=. -benchtime=3s .  # run benchmarks
```

## Notes

- `StructSliceToMap` and `GetFieldArray` use reflection. Make sure the field name and type parameter match exactly.
- `Optional.MustGet()` panics when the value is absent. Prefer `Get()` or `GetOrElse()` in production code.
- None of the types are safe for concurrent use without external synchronization.

## Author

[@liu-houliang](https://github.com/liu-houliang)