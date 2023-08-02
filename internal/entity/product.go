package entity

import (
	"errors"
	"time"

	"github.com/Nimbo1999/go-apis-go-expert/pkg/entity"
)

var (
	ErrIdIsRequired    = errors.New("id is required")
	ErrInvalidID       = errors.New("inalid id")
	ErrNameIsRequired  = errors.New("name is required")
	ErrPriceIsRequired = errors.New("price is required")
	ErrInvalidPrice    = errors.New("price is required")
)

type Product struct {
	ID        entity.ID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func NewProduct(name string, price float64) (*Product, error) {
	product := &Product{
		ID:        entity.NewID(),
		Price:     price,
		Name:      name,
		CreatedAt: time.Now(),
	}
	if err := product.Validate(); err != nil {
		return nil, err
	}
	return product, nil
}

func (p *Product) Validate() error {
	if p.ID.String() == "" {
		return ErrIdIsRequired
	}
	if _, err := entity.ParseID(p.ID.String()); err != nil {
		return ErrInvalidID
	}
	if p.Name == "" {
		return ErrNameIsRequired
	}
	if p.Price == 0.0 {
		return ErrPriceIsRequired
	}
	if p.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}
