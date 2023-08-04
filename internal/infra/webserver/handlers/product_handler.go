package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Nimbo1999/go-apis-go-expert/internal/dto"
	"github.com/Nimbo1999/go-apis-go-expert/internal/entity"
	"github.com/Nimbo1999/go-apis-go-expert/internal/infra/database"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

func (handler *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Printf("new request received: %s\n", r.Method)
	var createProductInput dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&createProductInput)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := entity.NewProduct(createProductInput.Name, createProductInput.Price)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = handler.ProductDB.Create(product); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
