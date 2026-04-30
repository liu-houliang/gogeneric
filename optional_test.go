package gogeneric_test

// Tests for optional.go
// Run: go test -run TestOptional .
// Run: go test -v -run TestOptional .

import (
	"fmt"
	"testing"

	"github.com/liu-houliang/gogeneric/v2"
)

func TestOptional_Some(t *testing.T) {
	opt := gogeneric.Some(42)
	if !opt.IsPresent() {
		t.Fatal("IsPresent = false, want true")
	}
	if opt.IsEmpty() {
		t.Fatal("IsEmpty = true, want false")
	}
	if v, ok := opt.Get(); !ok || v != 42 {
		t.Fatalf("Get = (%d, %v), want (42, true)", v, ok)
	}
}

func TestOptional_None(t *testing.T) {
	opt := gogeneric.None[int]()
	if opt.IsPresent() {
		t.Fatal("IsPresent = true, want false")
	}
	if !opt.IsEmpty() {
		t.Fatal("IsEmpty = false, want true")
	}
	v, ok := opt.Get()
	if ok || v != 0 {
		t.Fatalf("Get on None = (%d, %v), want (0, false)", v, ok)
	}
}

func TestOptional_GetOrElse(t *testing.T) {
	if v := gogeneric.Some(1).GetOrElse(99); v != 1 {
		t.Fatalf("GetOrElse(Some) = %d, want 1", v)
	}
	if v := gogeneric.None[int]().GetOrElse(99); v != 99 {
		t.Fatalf("GetOrElse(None) = %d, want 99", v)
	}
}

func TestOptional_GetOrElseGet(t *testing.T) {
	called := false
	supplier := func() int {
		called = true
		return 99
	}
	if v := gogeneric.Some(1).GetOrElseGet(supplier); v != 1 {
		t.Fatalf("GetOrElseGet(Some) = %d, want 1", v)
	}
	if called {
		t.Fatal("supplier should not be called when value is present")
	}
	if v := gogeneric.None[int]().GetOrElseGet(supplier); v != 99 {
		t.Fatalf("GetOrElseGet(None) = %d, want 99", v)
	}
	if !called {
		t.Fatal("supplier should be called when value is absent")
	}
}

func TestOptional_GetOrError(t *testing.T) {
	v, err := gogeneric.Some("hello").GetOrError("not found")
	if err != nil || v != "hello" {
		t.Fatalf("GetOrError(Some) = (%q, %v), want (hello, nil)", v, err)
	}
	_, err = gogeneric.None[string]().GetOrError("not found")
	if err == nil {
		t.Fatal("GetOrError(None) should return an error")
	}
}

func TestOptional_IfPresent(t *testing.T) {
	called := false
	gogeneric.Some(1).IfPresent(func(v int) { called = true })
	if !called {
		t.Fatal("IfPresent should call fn when value is present")
	}
	called = false
	gogeneric.None[int]().IfPresent(func(v int) { called = true })
	if called {
		t.Fatal("IfPresent should NOT call fn when value is absent")
	}
}

func TestOptional_MustGet_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("MustGet on None should panic")
		}
	}()
	gogeneric.None[int]().MustGet()
}

func TestOptional_MustGet_Value(t *testing.T) {
	if v := gogeneric.Some(42).MustGet(); v != 42 {
		t.Fatalf("MustGet = %d, want 42", v)
	}
}

func TestMapOptional(t *testing.T) {
	length := gogeneric.MapOptional(gogeneric.Some("hello"), func(s string) int { return len(s) })
	if v, _ := length.Get(); v != 5 {
		t.Fatalf("MapOptional = %d, want 5", v)
	}
	empty := gogeneric.MapOptional(gogeneric.None[string](), func(s string) int { return len(s) })
	if empty.IsPresent() {
		t.Fatal("MapOptional of None should be None")
	}
}

func TestFlatMapOptional(t *testing.T) {
	double := func(n int) gogeneric.Optional[int] {
		if n < 0 {
			return gogeneric.None[int]()
		}
		return gogeneric.Some(n * 2)
	}
	if v, _ := gogeneric.FlatMapOptional(gogeneric.Some(5), double).Get(); v != 10 {
		t.Fatalf("FlatMapOptional = %d, want 10", v)
	}
	if gogeneric.FlatMapOptional(gogeneric.Some(-1), double).IsPresent() {
		t.Fatal("FlatMapOptional → None when mapper returns None")
	}
	if gogeneric.FlatMapOptional(gogeneric.None[int](), double).IsPresent() {
		t.Fatal("FlatMapOptional of None should be None")
	}
}

func ExampleSome() {
	opt := gogeneric.Some("hello")
	fmt.Println(opt.GetOrElse("world"))
	// Output: hello
}

func ExampleNone() {
	opt := gogeneric.None[string]()
	fmt.Println(opt.GetOrElse("world"))
	// Output: world
}

func ExampleMapOptional() {
	length := gogeneric.MapOptional(gogeneric.Some("hello"), func(s string) int { return len(s) })
	fmt.Println(length.GetOrElse(0))
	// Output: 5
}
