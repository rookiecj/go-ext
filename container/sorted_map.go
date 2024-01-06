package container

import (
	"cmp"
	"slices"
)

// sort.Interface
// Len() int
// Less(i, j int) bool
// Swap(i, j int)

type SortedMap[K cmp.Ordered, V any] map[K]V

func (c *SortedMap[K, V]) SortedKeys(cmp func(a, b K) int) []K {
	var keys []K
	for key := range *c {
		keys = append(keys, key)
	}
	if cmp == nil {
		slices.Sort(keys)
	} else {
		slices.SortFunc(keys, cmp)
	}
	return keys
}
