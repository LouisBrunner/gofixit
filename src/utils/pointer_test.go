package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Pointerize_int(t *testing.T) {
	i := 5
	assert.Equal(t, &i, Pointerize(i))
}
