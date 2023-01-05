package blia

import "sort"

func IndexMap[T comparable](collection []T) map[T]int {
	m := make(map[T]int)
	for i, v := range collection {
		m[v] = i
	}
	return m
}

func SortAs[T comparable, U any](collection []U, indexes []T, keyFn func(U) T) {
	indexesMap := IndexMap(indexes)
	sort.Slice(collection, func(i, j int) bool {
		return indexesMap[keyFn(collection[i])] < indexesMap[keyFn(collection[j])]
	})
}
