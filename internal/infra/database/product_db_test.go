package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Nimbo1999/go-apis-go-expert/internal/entity"
	pkgEntity "github.com/Nimbo1999/go-apis-go-expert/pkg/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test_ProductCreate(t *testing.T) {
	productDB, err := createMemoryDB()
	assert.NoError(t, err)
	product, err := entity.NewProduct("Product 1", 10.)
	assert.NoError(t, err)
	err = productDB.Create(product)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.ID)
	assert.Equal(t, product.Name, "Product 1")
	assert.Equal(t, product.Price, 10.0)
	assert.NotNil(t, product.CreatedAt)
}

func Test_ProductFindAll(t *testing.T) {
	productDB, err := createMemoryDB()
	assert.NoError(t, err)
	for i := 0; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		assert.NoError(t, err)
		productDB.Create(product)
	}
	products, err := productDB.FindAll(1, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 0", products[0].Name)
	assert.Equal(t, "Product 9", products[9].Name)

	products, err = productDB.FindAll(2, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 10", products[0].Name)
	assert.Equal(t, "Product 19", products[9].Name)

	products, err = productDB.FindAll(3, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 4)
	assert.Equal(t, "Product 20", products[0].Name)
	assert.Equal(t, "Product 23", products[3].Name)

	products, err = productDB.FindAll(1, 10, "desc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 23", products[0].Name)
}

func Test_ProductFindById(t *testing.T) {
	productDB, err := createMemoryDB()
	assert.NoError(t, err)
	p, err := entity.NewProduct("Product 1", 10.)
	assert.NoError(t, err)
	err = productDB.Create(p)
	assert.NoError(t, err)
	product, err := productDB.FindById(p.ID.String())
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, p.ID, product.ID)
}

func Test_ProductFindById_ShouldNotFindProductThatDoesNotExist(t *testing.T) {
	productDB, err := createMemoryDB()
	assert.NoError(t, err)
	product, err := productDB.FindById(pkgEntity.NewID().String())
	assert.NotNil(t, err)
	assert.Nil(t, product)
}

func Test_ProductUpdate(t *testing.T) {
	productDB, err := createMemoryDB()
	assert.NoError(t, err)
	p, err := entity.NewProduct("Product 1", 10.)
	assert.NoError(t, err)
	err = productDB.Create(p)
	assert.NoError(t, err)
	p.Name = "Updated product"
	p.Price = 5.
	err = productDB.Update(p)
	assert.Nil(t, err)
	assert.Equal(t, p.Name, "Updated product")
	assert.Equal(t, p.Price, 5.0)
	assert.NotEmpty(t, p.ID)
}

func Test_ProductDelete(t *testing.T) {
	productDB, err := createMemoryDB()
	assert.NoError(t, err)

	p, err := entity.NewProduct("Product 1", 10.)
	assert.NoError(t, err)

	err = productDB.Create(p)
	assert.NoError(t, err)

	err = productDB.Delete(p.ID.String())
	assert.NoError(t, err)

	product, err := productDB.FindById(p.ID.String())
	assert.Error(t, err)
	assert.Nil(t, product)
}

func createMemoryDB() (*Product, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&entity.Product{})
	return NewProduct(db), err
}
