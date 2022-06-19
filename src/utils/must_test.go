package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Must_time_Works(t *testing.T) {
	tryThis := func() (time.Time, error) {
		return time.Parse("02/01/2006", "18/05/1991")
	}

	res, err := tryThis()
	if err != nil {
		t.Fatalf("failed to do something that should not fail")
	}

	assert.Equal(t, res, Must(tryThis()))
}

func Test_Must_time_Fails(t *testing.T) {
	tryThis := func() (time.Time, error) {
		return time.Parse("02/01/02", "18/05/1991")
	}

	assert.Panics(t, func() {
		Must(tryThis())
	})
}
