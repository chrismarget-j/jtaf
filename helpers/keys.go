package helpers

import (
	"golang.org/x/exp/constraints"
	"sort"
)

// Keys returns a slice of map keys
func Keys[A comparable, B interface{}](m map[A]B) []A {
	result := make([]A, len(m))
	i := 0
	for k := range m {
		result[i] = k
		i++
	}

	return result
}

// OrderedKeys returns a sorted slice of map keys
func OrderedKeys[A constraints.Ordered, B interface{}](m map[A]B) []A {
	result := Keys(m)

	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	return result
}