package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNullFloat64WhenNil(t *testing.T) {
	var id *float64 = nil
	subject := newNullFloat64(id)
	assert.NotNil(t, subject)
	assert.False(t, subject.Valid)
	assert.Zero(t, subject.Float64)
}

func TestNewNullFloat64WhenNotNil(t *testing.T) {
	val := float64(1.0)
	var id *float64 = &val
	subject := newNullFloat64(id)
	assert.NotNil(t, subject)
	assert.True(t, subject.Valid)
	assert.Equal(t, 1.0, subject.Float64)
}
