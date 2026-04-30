package gogeneric

import "errors"

// Optional represents a value that may or may not be present.
// It is the Go equivalent of Haskell's Maybe / Rust's Option<T>.
// The zero value of Optional is equivalent to None.
//
// Optional is not safe for concurrent use without external synchronization.
type Optional[T any] struct {
	value   T
	present bool
}

// Some wraps value in an Optional that is considered present.
//
// Example:
//
//	opt := Some(42)
//	opt.IsPresent() // true
//	opt.MustGet()   // 42
func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, present: true}
}

// None returns an empty Optional[T] with no value present.
//
// Example:
//
//	opt := None[string]()
//	opt.IsPresent() // false
func None[T any]() Optional[T] {
	return Optional[T]{}
}

// IsPresent reports whether a value is present.
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// IsEmpty reports whether no value is present.
func (o Optional[T]) IsEmpty() bool {
	return !o.present
}

// Get returns the value and true if present, or the zero value and false if empty.
//
// Example:
//
//	val, ok := opt.Get()
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.present
}

// MustGet returns the value or panics if no value is present.
//
// Use this only when you are certain the Optional is non-empty.
func (o Optional[T]) MustGet() T {
	if !o.present {
		panic("gogeneric: MustGet called on empty Optional")
	}
	return o.value
}

// GetOrElse returns the value if present, or fallback otherwise.
//
// Example:
//
//	name := opt.GetOrElse("anonymous")
func (o Optional[T]) GetOrElse(fallback T) T {
	if o.present {
		return o.value
	}
	return fallback
}

// GetOrElseGet returns the value if present, or the result of calling supplier.
// Use this to lazily compute the fallback value.
//
// Example:
//
//	name := opt.GetOrElseGet(func() string { return expensiveDefault() })
func (o Optional[T]) GetOrElseGet(supplier func() T) T {
	if o.present {
		return o.value
	}
	return supplier()
}

// GetOrError returns the value if present, or an error wrapping msg.
//
// Example:
//
//	val, err := opt.GetOrError("user not found")
func (o Optional[T]) GetOrError(msg string) (T, error) {
	if o.present {
		return o.value, nil
	}
	var zero T
	return zero, errors.New(msg)
}

// IfPresent calls fn with the value if present; otherwise it is a no-op.
//
// Example:
//
//	opt.IfPresent(func(v int) { fmt.Println(v) })
func (o Optional[T]) IfPresent(fn func(T)) {
	if o.present {
		fn(o.value)
	}
}

// MapOptional transforms the value inside an Optional using mapper.
// If the Optional is empty, None is returned.
//
// This is a package-level function (not a method) because Go does not
// support parameterized methods with additional type parameters.
//
// Example:
//
//	length := MapOptional(Some("hello"), func(s string) int { return len(s) })
//	// length == Some(5)
func MapOptional[T, R any](o Optional[T], mapper func(T) R) Optional[R] {
	if !o.present {
		return None[R]()
	}
	return Some(mapper(o.value))
}

// FlatMapOptional transforms the value using mapper which itself returns an Optional.
// If the input is empty, None is returned.
//
// Example:
//
//	parseAge := func(s string) Optional[int] { ... }
//	age := FlatMapOptional(Some("25"), parseAge)
func FlatMapOptional[T, R any](o Optional[T], mapper func(T) Optional[R]) Optional[R] {
	if !o.present {
		return None[R]()
	}
	return mapper(o.value)
}
