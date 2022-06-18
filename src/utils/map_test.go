package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SortedMap_String(t *testing.T) {
	in := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	expected := OrderedMap[string, int]{
		MapEntry[string, int]{
			Key:   "one",
			Value: 1,
		},
		MapEntry[string, int]{
			Key:   "three",
			Value: 3,
		},
		MapEntry[string, int]{
			Key:   "two",
			Value: 2,
		},
	}
	out := SortedMap(in)
	assert.Equal(t, expected, out)
}
