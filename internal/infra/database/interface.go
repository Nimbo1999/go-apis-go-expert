package database

import "github.com/Nimbo1999/go-apis-go-expert/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}
