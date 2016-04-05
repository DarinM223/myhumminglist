package main

import (
	"reflect"
	"testing"
)

func TestLRUMap_BasicOrderTest(t *testing.T) {
	lruMap := NewLRUMap()
	lruMap.Add(1, 2)
	lruMap.Add(2, 3)
	lruMap.Add(1, 4)

	expectedValues := []int{2, 1}
	if !reflect.DeepEqual(lruMap.Keys(), expectedValues) {
		t.Errorf("TestLRUMap_BasicOrderTest failed: want %v got %v", expectedValues, lruMap.Keys())
	}
}

func TestLRUMap_GetDoesntAffectLRU(t *testing.T) {
	lruMap := NewLRUMap()
	lruMap.Add(1, 2)
	lruMap.Add(2, 3)
	_, _ = lruMap.Get(1)

	expectedValues := []int{1, 2}
	if !reflect.DeepEqual(lruMap.Keys(), expectedValues) {
		t.Errorf("TestLRUMap_GetDoesntAffectLRU failed: want %v got %v", expectedValues, lruMap.Keys())
	}
}

func TestLRUMap_DeleteRemovesFromOrder(t *testing.T) {
	lruMap := NewLRUMap()
	lruMap.Add(1, 2)
	lruMap.Add(2, 3)
	lruMap.Add(3, 4)
	lruMap.Remove(2)

	expectedValues := []int{1, 3}
	if !reflect.DeepEqual(lruMap.Keys(), expectedValues) {
		t.Errorf("TestLRUMap_DeleteRemovesFromOrder failed: want %v got %v", expectedValues, lruMap.Keys())
	}
}
