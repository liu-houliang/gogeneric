package gogeneric

// Set is a generic, unordered collection of unique elements.
// It is backed by a map[T]struct{} and is not safe for concurrent use.
type Set[T comparable] struct {
	m map[T]struct{}
}

// NewSet creates a new Set optionally pre-populated with items.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	s.Contains(2) // true
func NewSet[T comparable](items ...T) Set[T] {
	s := Set[T]{m: make(map[T]struct{}, len(items))}
	for _, item := range items {
		s.m[item] = struct{}{}
	}
	return s
}

// Add inserts item into the set.
func (s *Set[T]) Add(item T) {
	s.m[item] = struct{}{}
}

// Remove deletes item from the set. No-op if not present.
func (s *Set[T]) Remove(item T) {
	delete(s.m, item)
}

// Contains reports whether item is in the set.
func (s Set[T]) Contains(item T) bool {
	_, ok := s.m[item]
	return ok
}

// Len returns the number of elements in the set.
func (s Set[T]) Len() int {
	return len(s.m)
}

// Slice returns all elements as an unordered slice.
func (s Set[T]) Slice() []T {
	result := make([]T, 0, len(s.m))
	for k := range s.m {
		result = append(result, k)
	}
	return result
}

// Union returns a new Set containing all elements from s and other.
//
// Example:
//
//	a := NewSet(1, 2, 3)
//	b := NewSet(3, 4, 5)
//	a.Union(b).Slice() // [1, 2, 3, 4, 5] (order not guaranteed)
func (s Set[T]) Union(other Set[T]) Set[T] {
	result := NewSet[T]()
	for k := range s.m {
		result.m[k] = struct{}{}
	}
	for k := range other.m {
		result.m[k] = struct{}{}
	}
	return result
}

// Intersection returns a new Set containing only elements present in both s and other.
//
// Example:
//
//	a := NewSet(1, 2, 3)
//	b := NewSet(2, 3, 4)
//	a.Intersection(b).Slice() // [2, 3] (order not guaranteed)
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	result := NewSet[T]()
	for k := range s.m {
		if _, ok := other.m[k]; ok {
			result.m[k] = struct{}{}
		}
	}
	return result
}

// Difference returns a new Set with elements in s but not in other.
//
// Example:
//
//	a := NewSet(1, 2, 3)
//	b := NewSet(2, 3, 4)
//	a.Difference(b).Slice() // [1]
func (s Set[T]) Difference(other Set[T]) Set[T] {
	result := NewSet[T]()
	for k := range s.m {
		if _, ok := other.m[k]; !ok {
			result.m[k] = struct{}{}
		}
	}
	return result
}

// IsSubset reports whether every element of s is also in other.
func (s Set[T]) IsSubset(other Set[T]) bool {
	for k := range s.m {
		if _, ok := other.m[k]; !ok {
			return false
		}
	}
	return true
}
