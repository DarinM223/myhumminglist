package main

import (
	"reflect"
	"testing"
)

var sortedInsertTests = []struct {
	array    []int
	value    int
	expected []int
}{
	{[]int{1, 2, 3, 4, 5, 6, 7}, 3, []int{1, 2, 3, 3, 4, 5, 6, 7}},
	{[]int{1, 2, 4, 5, 6, 7}, 3, []int{1, 2, 3, 4, 5, 6, 7}},
	{[]int{1, 2, 3, 4, 5, 6, 7}, 8, []int{1, 2, 3, 4, 5, 6, 7, 8}},
	{[]int{1, 2, 3, 4, 5, 6, 7}, 1, []int{1, 1, 2, 3, 4, 5, 6, 7}},
	{[]int{1, 3, 5, 7, 9, 11}, 6, []int{1, 3, 5, 6, 7, 9, 11}},
	{[]int{}, 1, []int{1}},
	{[]int{1}, 2, []int{1, 2}},
	{[]int{1}, 0, []int{0, 1}},
}

var sortedRemoveTests = []struct {
	array    []int
	value    int
	expected []int
}{
	{[]int{1, 2, 3, 4, 5, 6, 7}, 8, []int{1, 2, 3, 4, 5, 6, 7}},
	{[]int{1, 2, 3, 4, 5, 6, 7}, 3, []int{1, 2, 4, 5, 6, 7}},
	{[]int{}, 3, []int{}},
	{[]int{2, 3}, 3, []int{2}},
	{[]int{3}, 3, []int{}},
	{[]int{0, 3}, 0, []int{3}},
}

func TestSortedInsert(t *testing.T) {
	for _, test := range sortedInsertTests {
		result := sortedInsert(test.array, test.value)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("TestSortedInsert failed: want %v got %v", test.expected, result)
		}
	}
}

func TestSortedRemove(t *testing.T) {
	for _, test := range sortedRemoveTests {
		result := sortedRemove(test.array, test.value)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("TestSortedRemove failed: want %v got %v", test.expected, result)
		}
	}
}
