package database

import (
	"testing"

	"github.com/Nimbo1999/go-apis-go-expert/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.User{})
	user, _ := entity.NewUser("Matheus", "matheus@abc.com", "123456")
	userDb := NewUser(db)
	err = userDb.Create(user)
	assert.Nil(t, err)

	var userFound entity.User
	err = db.Find(&userFound, "id = ?", user.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, userFound.ID, user.ID)
	assert.Equal(t, userFound.Email, user.Email)
	assert.Equal(t, userFound.Name, user.Name)
	assert.True(t, userFound.ValidatePassword("123456"))
}

func TestFindByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.User{})
	user, _ := entity.NewUser("Matheus", "matheus@abc.com", "123456")
	userDb := NewUser(db)
	err = userDb.Create(user)
	assert.Nil(t, err)

	userFound, err := userDb.FindByEmail("matheus@abc.com")
	assert.Nil(t, err)
	assert.NotNil(t, userFound)
	assert.Equal(t, userFound.ID, user.ID)
	assert.Equal(t, userFound.Name, user.Name)
	assert.Equal(t, userFound.Email, user.Email)
	assert.True(t, userFound.ValidatePassword("123456"))
	assert.NotEqual(t, userFound.Password, "123456")
}
