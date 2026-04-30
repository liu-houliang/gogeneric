package gogeneric_test

// Tests for generic.go
// Run: go test -run TestFilter .
// Run: go test -run TestMap .
// Run: go test -v .

import (
	"fmt"
	"testing"

	"github.com/liu-houliang/gogeneric/v2"
)

// ---------------------------------------------------------------------------
// Deprecated helpers (kept for backward compatibility)
// ---------------------------------------------------------------------------

func TestMin(t *testing.T) {
	if got := gogeneric.Min(336, 445); got != 336 {
		t.Fatalf("Min(336, 445) = %d, want 336", got)
	}
	if got := gogeneric.Min(3.14, 2.72); got != 2.72 {
		t.Fatalf("Min(3.14, 2.72) = %v, want 2.72", got)
	}
}

func TestMax(t *testing.T) {
	if got := gogeneric.Max(336, 445); got != 445 {
		t.Fatalf("Max(336, 445) = %d, want 445", got)
	}
}

func TestCond(t *testing.T) {
	if got := gogeneric.Cond(78 > 60, "pass", "fail"); got != "pass" {
		t.Fatalf("Cond = %q, want \"pass\"", got)
	}
	if got := gogeneric.Cond(30 > 60, "pass", "fail"); got != "fail" {
		t.Fatalf("Cond = %q, want \"fail\"", got)
	}
}

func ExampleCond() {
	result := gogeneric.Cond(10 > 5, "big", "small")
	fmt.Println(result)
	// Output: big
}

func TestSortSlice(t *testing.T) {
	list := []int{3, 1, 2, 4, 5}
	gogeneric.SortSlice(list, func(i, j int) bool { return i < j })
	want := []int{1, 2, 3, 4, 5}
	for i, v := range list {
		if v != want[i] {
			t.Fatalf("SortSlice[%d] = %d, want %d", i, v, want[i])
		}
	}
}

func TestSliceIndex(t *testing.T) {
	s := []string{"alice", "bob", "carol"}
	if got := gogeneric.SliceIndex(s, "bob"); got != 1 {
		t.Fatalf("SliceIndex = %d, want 1", got)
	}
	if got := gogeneric.SliceIndex(s, "dave"); got != -1 {
		t.Fatalf("SliceIndex = %d, want -1", got)
	}
}

func TestCompareSlice(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{1, 2, 3}
	c := []int{1, 2, 4}
	if !gogeneric.CompareSlice(a, b) {
		t.Fatal("CompareSlice(equal) = false, want true")
	}
	if gogeneric.CompareSlice(a, c) {
		t.Fatal("CompareSlice(different) = true, want false")
	}
	if gogeneric.CompareSlice(a, []int{1, 2}) {
		t.Fatal("CompareSlice(different length) = true, want false")
	}
}

func ExampleCompareSlice() {
	fmt.Println(gogeneric.CompareSlice([]int{1, 2, 3}, []int{1, 2, 3}))
	fmt.Println(gogeneric.CompareSlice([]int{1, 2, 3}, []int{1, 2, 4}))
	// Output:
	// true
	// false
}

func TestMapKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := gogeneric.MapKeys(m)
	if len(keys) != 3 {
		t.Fatalf("MapKeys len = %d, want 3", len(keys))
	}
}

func TestMapToSlice(t *testing.T) {
	m := map[string]int{"x": 10}
	vals := gogeneric.MapToSlice(m)
	if len(vals) != 1 || vals[0] != 10 {
		t.Fatalf("MapToSlice = %v, want [10]", vals)
	}
}

func TestSliceToMap(t *testing.T) {
	m := gogeneric.SliceToMap([]int{1, 2, 3})
	for _, v := range []int{1, 2, 3} {
		if m[v] != v {
			t.Fatalf("SliceToMap[%d] = %d, want %d", v, m[v], v)
		}
	}
}

// ---------------------------------------------------------------------------
// Struct utilities
// ---------------------------------------------------------------------------

type student struct {
	Name string
	Age  uint
}

func TestStructSliceToMap(t *testing.T) {
	students := []student{{"alice", 20}, {"bob", 25}}
	m := gogeneric.StructSliceToMap[string]("Name", students)
	if len(m) != 2 {
		t.Fatalf("StructSliceToMap len = %d, want 2", len(m))
	}
	if m["alice"].Age != 20 {
		t.Fatalf("StructSliceToMap[alice].Age = %d, want 20", m["alice"].Age)
	}
	// wrong key type → empty map
	empty := gogeneric.StructSliceToMap[int]("Name", students)
	if len(empty) != 0 {
		t.Fatal("StructSliceToMap with wrong key type should return empty map")
	}
}

func TestGetFieldArray(t *testing.T) {
	students := []student{{"alice", 20}, {"bob", 25}}
	names := gogeneric.GetFieldArray[string]("Name", students)
	if len(names) != 2 || names[0] != "alice" || names[1] != "bob" {
		t.Fatalf("GetFieldArray = %v, want [alice bob]", names)
	}
	empty := gogeneric.GetFieldArray[int]("Name", students)
	if len(empty) != 0 {
		t.Fatal("GetFieldArray with wrong type should return empty slice")
	}
}

// ---------------------------------------------------------------------------
// Functional helpers
// ---------------------------------------------------------------------------

func TestFilter(t *testing.T) {
	evens := gogeneric.Filter([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 0 })
	if len(evens) != 2 || evens[0] != 2 || evens[1] != 4 {
		t.Fatalf("Filter = %v, want [2 4]", evens)
	}
	// empty input
	if got := gogeneric.Filter([]int{}, func(n int) bool { return true }); len(got) != 0 {
		t.Fatalf("Filter(empty) = %v, want []", got)
	}
}

func ExampleFilter() {
	evens := gogeneric.Filter([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 0 })
	fmt.Println(evens)
	// Output: [2 4]
}

func TestMap(t *testing.T) {
	doubled := gogeneric.Map([]int{1, 2, 3}, func(n int) int { return n * 2 })
	want := []int{2, 4, 6}
	for i, v := range doubled {
		if v != want[i] {
			t.Fatalf("Map[%d] = %d, want %d", i, v, want[i])
		}
	}
}

func ExampleMap() {
	strs := gogeneric.Map([]int{1, 2, 3}, func(n int) string { return fmt.Sprintf("%d", n*n) })
	fmt.Println(strs)
	// Output: [1 4 9]
}

func TestReduce(t *testing.T) {
	sum := gogeneric.Reduce([]int{1, 2, 3, 4, 5}, 0, func(acc, n int) int { return acc + n })
	if sum != 15 {
		t.Fatalf("Reduce sum = %d, want 15", sum)
	}
	// empty slice returns initial
	if got := gogeneric.Reduce([]int{}, 42, func(acc, n int) int { return acc + n }); got != 42 {
		t.Fatalf("Reduce(empty) = %d, want 42", got)
	}
}

func ExampleReduce() {
	sum := gogeneric.Reduce([]int{1, 2, 3, 4, 5}, 0, func(acc, n int) int { return acc + n })
	fmt.Println(sum)
	// Output: 15
}

func TestForEach(t *testing.T) {
	var got []int
	gogeneric.ForEach([]int{1, 2, 3}, func(n int) { got = append(got, n*2) })
	want := []int{2, 4, 6}
	for i, v := range got {
		if v != want[i] {
			t.Fatalf("ForEach result[%d] = %d, want %d", i, v, want[i])
		}
	}
}

func TestChunk(t *testing.T) {
	chunks := gogeneric.Chunk([]int{1, 2, 3, 4, 5}, 2)
	if len(chunks) != 3 {
		t.Fatalf("Chunk len = %d, want 3", len(chunks))
	}
	if len(chunks[2]) != 1 || chunks[2][0] != 5 {
		t.Fatalf("Chunk last = %v, want [5]", chunks[2])
	}
	if gogeneric.Chunk([]int{1, 2}, 0) != nil {
		t.Fatal("Chunk(size=0) should return nil")
	}
}

func ExampleChunk() {
	chunks := gogeneric.Chunk([]int{1, 2, 3, 4, 5}, 2)
	fmt.Println(chunks)
	// Output: [[1 2] [3 4] [5]]
}

func TestFlatten(t *testing.T) {
	flat := gogeneric.Flatten([][]int{{1, 2}, {3, 4}, {5}})
	want := []int{1, 2, 3, 4, 5}
	for i, v := range flat {
		if v != want[i] {
			t.Fatalf("Flatten[%d] = %d, want %d", i, v, want[i])
		}
	}
}

func ExampleFlatten() {
	flat := gogeneric.Flatten([][]int{{1, 2}, {3}, {4, 5}})
	fmt.Println(flat)
	// Output: [1 2 3 4 5]
}

func TestUnique(t *testing.T) {
	u := gogeneric.Unique([]int{1, 2, 2, 3, 1, 4})
	if len(u) != 4 {
		t.Fatalf("Unique len = %d, want 4", len(u))
	}
	// order preserved
	if u[0] != 1 || u[1] != 2 || u[2] != 3 || u[3] != 4 {
		t.Fatalf("Unique = %v, want [1 2 3 4]", u)
	}
}

func ExampleUnique() {
	u := gogeneric.Unique([]int{1, 2, 2, 3, 1})
	fmt.Println(u)
	// Output: [1 2 3]
}

// ---------------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------------

func BenchmarkFilter(b *testing.B) {
	s := make([]int, 1000)
	for i := range s {
		s[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gogeneric.Filter(s, func(n int) bool { return n%2 == 0 })
	}
}

func BenchmarkMap(b *testing.B) {
	s := make([]int, 1000)
	for i := range s {
		s[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gogeneric.Map(s, func(n int) int { return n * 2 })
	}
}

func BenchmarkUnique(b *testing.B) {
	s := make([]int, 1000)
	for i := range s {
		s[i] = i % 100 // lots of duplicates
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gogeneric.Unique(s)
	}
}
