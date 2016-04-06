package main

import (
	"errors"
	"fmt"
)

// sortedInsert inserts a value into a sorted array of integers
func sortedInsert(arr []int, value int) []int {
	if len(arr) == 0 {
		return []int{value}
	} else if value < arr[0] {
		return append([]int{value}, arr...)
	} else if value > arr[len(arr)-1] {
		return append(arr, value)
	}

	start, end := 0, len(arr)
	index := 0
	for start <= end {
		mid := (start + end) / 2
		if mid == start || mid == end {
			index = mid
			break
		} else if value == arr[mid] {
			index = mid
			break
		} else if value < arr[mid] {
			end = mid
		} else {
			start = mid
		}
	}

	bottomHalf := make([]int, index+1)
	copy(bottomHalf, arr[:index+1])
	bottomHalf = append(bottomHalf, value)

	return append(bottomHalf, arr[index+1:]...)
}

// sortedRemove removes a value from a sorted array of integers
func sortedRemove(arr []int, value int) []int {
	if len(arr) == 0 {
		return arr
	}

	start, end := 0, len(arr)
	index := 0
	for start <= end {
		mid := (start + end) / 2
		if value == arr[mid] {
			index = mid
			break
		} else if mid == start {
			return arr
		} else if value < arr[mid] {
			end = mid
		} else {
			start = mid
		}
	}

	return append(arr[:index], arr[index+1:]...)
}

// OrderedMap is a map that allows you to iterate
// over the keys in numerical order
type OrderedMap struct {
	// sorted keys of indexes
	keys []int
	// map of index to value
	dict map[int]interface{}
}

// NewOrderedMap creates a new ordered map
func NewOrderedMap() OrderedMap {
	return OrderedMap{
		keys: []int{},
		dict: make(map[int]interface{}),
	}
}

// Add adds a new key-value mapping to the map
func (m *OrderedMap) Add(key int, value interface{}) {
	if _, ok := m.dict[key]; !ok {
		m.keys = sortedInsert(m.keys, key)
	}
	m.dict[key] = value
}

// Remove removes a key from the map
func (m *OrderedMap) Remove(key int) {
	delete(m.dict, key)
	m.keys = sortedRemove(m.keys, key)
}

// Get retrieves a value from an integer key
func (m *OrderedMap) Get(key int) (interface{}, error) {
	if value, ok := m.dict[key]; ok {
		return value, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Key %d is not in the OrderedMap", key))
	}
}

// Contains returns true if the map contains the key, false otherwise
func (m *OrderedMap) Contains(key int) bool {
	_, ok := m.dict[key]
	return ok
}
