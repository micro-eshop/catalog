package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductIdValidationWhenIsInvalid(t *testing.T) {
	var id ProductId = -1
	subject := ValidateProductId(id)
	assert.NotNil(t, subject)
	assert.EqualError(t, subject, "ProductId must be greater than 0")
}

func TestProductIdValidationWhenIsValid(t *testing.T) {
	var id ProductId = 2
	subject := ValidateProductId(id)
	assert.Nil(t, subject)

}
