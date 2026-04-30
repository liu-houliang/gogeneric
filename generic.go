package gogeneric

import (
	"reflect"
	"sort"
)

// Number constraint normal number type
type Number interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

// SliceFn implements sort.Interface for a slice of T.
type SliceFn[T any] struct {
	s    []T
	less func(T, T) bool
}

func (s SliceFn[T]) Len() int           { return len(s.s) }
func (s SliceFn[T]) Swap(i, j int)      { s.s[i], s.s[j] = s.s[j], s.s[i] }
func (s SliceFn[T]) Less(i, j int) bool { return s.less(s.s[i], s.s[j]) }

// Min returns the minimum of two values.
//
// Deprecated: Use the built-in min() function introduced in Go 1.21.
func Min[T Number](base, comp T) T {
	if base > comp {
		return comp
	}
	return base
}

// Max returns the maximum of two values.
//
// Deprecated: Use the built-in max() function introduced in Go 1.21.
func Max[T Number](base, comp T) T {
	if base < comp {
		return comp
	}
	return base
}

// Cond is a generic ternary operator.
// It returns a if condition is true, otherwise returns b.
//
// Example:
//
//	result := Cond(score > 60, "pass", "fail")
func Cond[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

// SortSlice sorts a slice using the provided less function.
//
// Deprecated: Use slices.SortFunc from the standard library (Go 1.21+).
func SortSlice[T any](slice []T, less func(T, T) bool) {
	sort.Sort(SliceFn[T]{slice, less})
}

// SliceIndex finds the first index of ele in slice.
// Returns -1 if not found.
//
// Deprecated: Use slices.Index from the standard library (Go 1.21+).
func SliceIndex[T comparable](slice []T, ele T) int {
	for k, v := range slice {
		if v == ele {
			return k
		}
	}
	return -1
}

// CompareSlice reports whether two slices are equal:
// same length and all elements equal in order.
//
// Deprecated: Use slices.Equal from the standard library (Go 1.21+).
func CompareSlice[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	if (a == nil) != (b == nil) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// MapKeys returns all keys of the map as a slice.
// The order of the returned slice is not guaranteed.
//
// Deprecated: Use maps.Keys from the standard library (Go 1.21+).
func MapKeys[Key comparable, Val any](m map[Key]Val) []Key {
	s := make([]Key, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

// MapToSlice returns all values of the map as a slice.
// The order of the returned slice is not guaranteed.
//
// Deprecated: Use maps.Values from the standard library (Go 1.21+).
func MapToSlice[Key comparable, Val any](m map[Key]Val) []Val {
	s := make([]Val, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

// SliceToMap converts a slice to a map[T]T using each element as both key and value.
func SliceToMap[T comparable](slice []T) map[T]T {
	m := make(map[T]T, len(slice))
	for _, v := range slice {
		m[v] = v
	}
	return m
}

// StructSliceToMap converts a struct slice to a map keyed by the given field name.
// The first occurrence wins if multiple elements share the same key value.
//
// key is the struct field name used as the map key.
// Key is the type of that field (must be comparable).
// Value is the struct type.
//
// Returns an empty map if slice is empty, the field does not exist,
// or the field type does not match Key.
func StructSliceToMap[Key comparable, Value any](key string, slice []Value) map[Key]Value {
	maps := make(map[Key]Value, len(slice))
	if len(slice) == 0 {
		return maps
	}
	typeT := reflect.TypeOf(slice[0])
	if typeT == nil || typeT.Kind() != reflect.Struct {
		return maps
	}
	field, ok := typeT.FieldByName(key)
	if !ok {
		return maps
	}
	var zeroKey Key
	if reflect.TypeOf(zeroKey) != field.Type {
		return maps
	}
	for _, v := range slice {
		rv := reflect.ValueOf(v)
		fv := rv.FieldByName(key)
		if !fv.IsValid() {
			continue
		}
		k := fv.Interface().(Key)
		if _, exists := maps[k]; !exists {
			maps[k] = v
		}
	}
	return maps
}

// GetFieldArray extracts a specific field from each element of a struct slice.
//
// fieldName is the struct field name to extract.
// Field is the type of that field.
// T is the struct type.
//
// Returns an empty slice if slice is empty, the field does not exist,
// or the field type does not match Field.
func GetFieldArray[Field any, T any](fieldName string, slice []T) []Field {
	fields := make([]Field, 0, len(slice))
	if len(slice) == 0 {
		return fields
	}
	typeT := reflect.TypeOf(slice[0])
	if typeT == nil || typeT.Kind() != reflect.Struct {
		return fields
	}
	sf, ok := typeT.FieldByName(fieldName)
	if !ok {
		return fields
	}
	var zeroField Field
	if reflect.TypeOf(zeroField) != sf.Type {
		return fields
	}
	for _, v := range slice {
		rv := reflect.ValueOf(v)
		fv := rv.FieldByName(fieldName)
		if !fv.IsValid() {
			continue
		}
		fields = append(fields, fv.Interface().(Field))
	}
	return fields
}

// ---------------------------------------------------------------------------
// Functional slice helpers — NOT available in the standard library
// ---------------------------------------------------------------------------

// Filter returns a new slice containing only the elements for which
// predicate returns true.
//
// Example:
//
//	evens := Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
//	// evens == []int{2, 4}
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map transforms each element of slice using the mapper function and returns
// a new slice of the results.
//
// Example:
//
//	doubled := Map([]int{1, 2, 3}, func(n int) int { return n * 2 })
//	// doubled == []int{2, 4, 6}
func Map[T, R any](slice []T, mapper func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = mapper(v)
	}
	return result
}

// Reduce reduces slice to a single value by applying reducer to each element
// from left to right, starting with initial.
//
// Example:
//
//	sum := Reduce([]int{1, 2, 3, 4}, 0, func(acc, n int) int { return acc + n })
//	// sum == 10
func Reduce[T, R any](slice []T, initial R, reducer func(R, T) R) R {
	acc := initial
	for _, v := range slice {
		acc = reducer(acc, v)
	}
	return acc
}

// ForEach calls fn on each element of slice. Unlike range, this is chainable
// and can be used in a functional pipeline.
//
// Example:
//
//	ForEach([]string{"a", "b"}, func(s string) { fmt.Println(s) })
func ForEach[T any](slice []T, fn func(T)) {
	for _, v := range slice {
		fn(v)
	}
}

// Chunk splits slice into sub-slices of at most size elements.
// The last chunk may be smaller than size if the slice length is not evenly divisible.
//
// Example:
//
//	chunks := Chunk([]int{1, 2, 3, 4, 5}, 2)
//	// chunks == [][]int{{1, 2}, {3, 4}, {5}}
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	n := (len(slice) + size - 1) / size
	result := make([][]T, 0, n)
	for size < len(slice) {
		result = append(result, slice[:size])
		slice = slice[size:]
	}
	if len(slice) > 0 {
		result = append(result, slice)
	}
	return result
}

// Flatten merges a slice of slices into a single flat slice.
//
// Example:
//
//	flat := Flatten([][]int{{1, 2}, {3, 4}, {5}})
//	// flat == []int{1, 2, 3, 4, 5}
func Flatten[T any](slices [][]T) []T {
	total := 0
	for _, s := range slices {
		total += len(s)
	}
	result := make([]T, 0, total)
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}

// Unique returns a new slice with duplicate elements removed,
// preserving the first occurrence order.
//
// Example:
//
//	u := Unique([]int{1, 2, 2, 3, 1})
//	// u == []int{1, 2, 3}
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
