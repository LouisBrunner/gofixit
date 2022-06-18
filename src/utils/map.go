package utils

import (
	"sort"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
)

type MapEntry[K constraints.Ordered, V any] struct {
	Key   K
	Value V
}

type OrderedMap[K constraints.Ordered, V any] []MapEntry[K, V]

func SortedMap[K constraints.Ordered, V any](in map[K]V) OrderedMap[K, V] {
	keys := maps.Keys(in)
	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	out := make([]MapEntry[K, V], 0, len(keys))
	for _, key := range keys {
		out = append(out, MapEntry[K, V]{
			Key:   key,
			Value: in[key],
		})
	}

	return out
}
