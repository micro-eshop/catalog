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

func TestProductIdsValidationWhenIsValid(t *testing.T) {
	var ids = []ProductId{ProductId(1), ProductId(2), ProductId(3)}
	subject := ValidateProductIds(ids)
	assert.Nil(t, subject)
}

func TestProductIdsValidationWhenIsInValid(t *testing.T) {
	var ids = []ProductId{ProductId(1), ProductId(-1), ProductId(3)}
	subject := ValidateProductIds(ids)
	assert.EqualError(t, subject, "ProductId must be greater than 0")
}
