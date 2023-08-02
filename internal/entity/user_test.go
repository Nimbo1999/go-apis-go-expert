package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("Matheus", "matheus@abc.com", "123456")
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.Equal(t, user.Name, "Matheus")
	assert.Equal(t, user.Email, "matheus@abc.com")
	assert.NotEmpty(t, user.Password)
	assert.NotEmpty(t, user.ID)
}

func TestUser_ValidatePassword(t *testing.T) {
	user, err := NewUser("Matheus", "matheus@abc.com", "123456")
	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword("123456"))
	assert.False(t, user.ValidatePassword("1234567"))
	assert.NotEqual(t, user.Password, "123456")
}
