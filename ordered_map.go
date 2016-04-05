package main

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

type orderedMap struct {
	// sorted keys of indexes
	keys []int
	// map of index to value
	dict map[int]interface{}
}

func newOrderedMap() orderedMap {
	return orderedMap{
		keys: []int{},
		dict: make(map[int]interface{}),
	}
}

func (m *orderedMap) add(key int, value interface{}) {
	if _, ok := m.dict[key]; !ok {
		m.keys = sortedInsert(m.keys, key)
	}
	m.dict[key] = value
}

func (m *orderedMap) remove(key int) {
	delete(m.dict, key)
	m.keys = sortedRemove(m.keys, key)
}
