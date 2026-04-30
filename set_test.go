package gogeneric_test

// Tests for set.go
// Run: go test -run TestSet .
// Run: go test -v -run TestSet .

import (
	"fmt"
	"testing"

	"github.com/liu-houliang/gogeneric/v2"
)

func TestSet_AddContains(t *testing.T) {
	s := gogeneric.NewSet(1, 2, 3)
	if !s.Contains(2) {
		t.Fatal("Contains(2) = false, want true")
	}
	if s.Contains(99) {
		t.Fatal("Contains(99) = true, want false")
	}
	s.Add(99)
	if !s.Contains(99) {
		t.Fatal("Contains(99) after Add = false, want true")
	}
}

func TestSet_Remove(t *testing.T) {
	s := gogeneric.NewSet("a", "b", "c")
	s.Remove("b")
	if s.Contains("b") {
		t.Fatal("Contains(b) after Remove = true, want false")
	}
	if s.Len() != 2 {
		t.Fatalf("Len after Remove = %d, want 2", s.Len())
	}
	// remove non-existent should be no-op
	s.Remove("zzz")
	if s.Len() != 2 {
		t.Fatalf("Len after Remove(non-existent) = %d, want 2", s.Len())
	}
}

func TestSet_Len(t *testing.T) {
	empty := gogeneric.NewSet[int]()
	if empty.Len() != 0 {
		t.Fatalf("empty set Len = %d, want 0", empty.Len())
	}
	s := gogeneric.NewSet(1, 1, 2) // duplicates in NewSet
	if s.Len() != 2 {
		t.Fatalf("NewSet with duplicates Len = %d, want 2", s.Len())
	}
}

func TestSet_Union(t *testing.T) {
	a := gogeneric.NewSet(1, 2, 3)
	b := gogeneric.NewSet(3, 4, 5)
	u := a.Union(b)
	if u.Len() != 5 {
		t.Fatalf("Union Len = %d, want 5", u.Len())
	}
	for _, v := range []int{1, 2, 3, 4, 5} {
		if !u.Contains(v) {
			t.Fatalf("Union missing %d", v)
		}
	}
}

func TestSet_Intersection(t *testing.T) {
	a := gogeneric.NewSet(1, 2, 3)
	b := gogeneric.NewSet(2, 3, 4)
	inter := a.Intersection(b)
	if inter.Len() != 2 {
		t.Fatalf("Intersection Len = %d, want 2", inter.Len())
	}
	if !inter.Contains(2) || !inter.Contains(3) {
		t.Fatalf("Intersection = %v, want {2,3}", inter.Slice())
	}
	// no overlap → empty
	c := gogeneric.NewSet(10, 11)
	if a.Intersection(c).Len() != 0 {
		t.Fatal("Intersection(no overlap) should be empty")
	}
}

func TestSet_Difference(t *testing.T) {
	a := gogeneric.NewSet(1, 2, 3)
	b := gogeneric.NewSet(2, 3, 4)
	diff := a.Difference(b)
	if diff.Len() != 1 || !diff.Contains(1) {
		t.Fatalf("Difference = %v, want {1}", diff.Slice())
	}
}

func TestSet_IsSubset(t *testing.T) {
	a := gogeneric.NewSet(1, 2)
	b := gogeneric.NewSet(1, 2, 3)
	if !a.IsSubset(b) {
		t.Fatal("IsSubset = false, want true")
	}
	if b.IsSubset(a) {
		t.Fatal("IsSubset (reversed) = true, want false")
	}
	// empty set is subset of everything
	empty := gogeneric.NewSet[int]()
	if !empty.IsSubset(a) {
		t.Fatal("empty set IsSubset = false, want true")
	}
}

func ExampleNewSet() {
	s := gogeneric.NewSet(1, 2, 3)
	s.Add(4)
	fmt.Println(s.Contains(4))
	fmt.Println(s.Contains(99))
	// Output:
	// true
	// false
}
