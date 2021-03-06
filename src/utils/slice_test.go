package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MapSlice_string2int(t *testing.T) {
	in := []string{
		"one",
		"two",
		"three",
	}
	expected := []int{
		3,
		3,
		5,
	}
	out := MapSlice(in, func(in string) int {
		return len(in)
	})
	assert.Equal(t, expected, out)
}
