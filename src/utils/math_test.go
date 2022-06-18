package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Min_int(t *testing.T) {
	tests := []struct {
		name   string
		first  int
		second int
		result int
	}{
		{
			name:   "greater",
			first:  42,
			second: -3,
			result: -3,
		},
		{
			name:   "lesser",
			first:  1,
			second: 3,
			result: 1,
		},
		{
			name:   "equal",
			first:  4,
			second: 4,
			result: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.result, Min(tt.first, tt.second))
		})
	}
}
