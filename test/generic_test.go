package test

import (
	"github.com/liu-houliang/gogeneric"
	"testing"
)

func TestMin(t *testing.T) {
	a := 336
	b := 445
	c := gogeneric.Min(a, b)
	// Min Number:  336
	t.Log("Min Number: ", c)
}

func TestMax(t *testing.T) {
	a := 336
	b := 445
	c := gogeneric.Max(a, b)
	//  Max Number:  445
	t.Log("Max Number: ", c)
}

func TestCond(t *testing.T) {
	score := 78
	// result := "wrong"
	// if score > 60 {
	// 	result = "right"
	// }
	result := gogeneric.Cond(score > 60, "right", "wrong")
	// Cond Result:  right
	t.Log("Cond Result: ", result)
}

func TestSortSlice(t *testing.T) {

	list := []int{3, 1, 2, 4, 5}
	gogeneric.SortSlice(list, func(i, j int) bool {
		return i < j
	})
	// SortSlice Result:  [1 2 3 4 5]
	t.Log("SortSlice Result: ", list)
}

func TestSliceIndex(t *testing.T) {

	students := []string{
		"wangmiao",
		"yewenjie",
	}
	res := gogeneric.SliceIndex(students, "luoji")
	// CompareSlice Result:  0
	t.Log("SliceIndex Result: ", res)
}

func TestCompareSlice(t *testing.T) {

	listA := []int{3, 1, 2, 4, 5}
	listB := []int{1, 2, 3, 4, 5}
	res := gogeneric.CompareSlice(listA, listB)
	// CompareSlice Result:  0
	t.Log("CompareSlice Result: ", res)
}

func TestMapKeys(t *testing.T) {

	mapTest := map[string]string{
		"a": "e",
		"b": "x",
		"c": "z",
	}
	list := gogeneric.MapKeys(mapTest)
	// MapKeys Result:  [b c a]
	t.Log("MapKeys Result: ", list)
}

func TestMapToSlice(t *testing.T) {

	mapTest := map[string]string{
		"a": "e",
		"b": "x",
		"c": "z",
	}
	list := gogeneric.MapToSlice(mapTest)
	// MapToSlice Result:  [e x z]
	t.Log("MapToSlice Result: ", list)
}

func TestSliceToMap(t *testing.T) {

	list := []int{3, 1, 2, 4, 5}
	res := gogeneric.SliceToMap(list)
	// SliceToMap Result:  map[1:1 2:2 3:3 4:4 5:5]
	t.Log("SliceToMap Result: ", res)
}

type Student struct {
	Name string
	Age  uint
}

func TestStructSliceToMap(t *testing.T) {

	students := []Student{
		{"wangmiao", 23},
		{"yewenjie", 65},
	}
	res := gogeneric.StructSliceToMap[string]("Name", students)
	// StructSliceToMap Result:  map[wangmiao:{wangmiao 23} yewenjie:{yewenjie 65}]
	t.Log("StructSliceToMap Result: ", res)
}

func TestGetFieldArray(t *testing.T) {

	students := []Student{
		{"wangmiao", 23},
		{"yewenjie", 65},
	}
	res := gogeneric.GetFieldArray[string]("Name", students)
	// GetFieldArray Result:  [wangmiao yewenjie]
	t.Log("GetFieldArray Result: ", res)
}
