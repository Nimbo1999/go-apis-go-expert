package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduc(t *testing.T) {
	p, err := NewProduct("Product 1", 10)
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.NotEmpty(t, p.ID)
	assert.Equal(t, p.Name, "Product 1")
	assert.Equal(t, p.Price, 10.)
}

func TestProductNameIsRequired(t *testing.T) {
	p, err := NewProduct("", 10.)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNameIsRequired, err)
	assert.EqualError(t, ErrNameIsRequired, err.Error())
	assert.Nil(t, p)
}

func TestProductPriceIsRequired(t *testing.T) {
	p, err := NewProduct("Product 1", 0.)
	assert.NotNil(t, err)
	assert.Nil(t, p)
	assert.Equal(t, ErrPriceIsRequired, err)
	assert.EqualError(t, ErrPriceIsRequired, err.Error())
}

func TestProductPriceIsInvalid(t *testing.T) {
	p, err := NewProduct("Product 1", -1.)
	assert.NotNil(t, err)
	assert.Nil(t, p)
	assert.Equal(t, ErrInvalidPrice, err)
	assert.EqualError(t, ErrInvalidPrice, err.Error())
}

func TestProductValidate(t *testing.T) {
	p, err := NewProduct("Product 1", 100)
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Nil(t, p.Validate())
}
