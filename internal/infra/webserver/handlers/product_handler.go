package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Nimbo1999/go-apis-go-expert/internal/dto"
	"github.com/Nimbo1999/go-apis-go-expert/internal/entity"
	"github.com/Nimbo1999/go-apis-go-expert/internal/infra/database"
	pkgEntity "github.com/Nimbo1999/go-apis-go-expert/pkg/entity"
	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// Create product godoc
// @Summary     Create product
// @Description Create product
// @Tags        products
// @Accept      json
// @Produce     json
// @Param       request       body        dto.CreateProductInput   true    "Product request payload"
// @Success     201           {object}    entity.Product
// @Failure     400
// @Failure     500           {object}    Error
// @Router      /product   [post]
// @Security    ApiKeyAuth
func (handler *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createProductInput dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&createProductInput)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	product, err := entity.NewProduct(createProductInput.Name, createProductInput.Price)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	if err = handler.ProductDB.Create(product); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
	}
}

// Get product godoc
// @Summary     Get product
// @Description Get product by id
// @Tags        products
// @Accept      json
// @Produce     json
// @Param       product_id    path        string                   true    "Product Id"
// @Success     200           {object}    entity.Product
// @Failure     400           {object}    Error
// @Failure     404           {object}    Error
// @Failure     500           {object}    Error
// @Router      /product/{product_id}     [get]
// @Security    ApiKeyAuth
func (handler *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Println("request does not have an id parameter.")
		w.WriteHeader(http.StatusBadRequest)
		errorObj := Error{Message: "you must provide an id to the url"}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	product, err := handler.ProductDB.FindById(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(&product); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
	}
}

// Update product godoc
// @Summary     Update products
// @Description Update product
// @Tags        products
// @Accept      json
// @Produce     json
// @Param       product_id    path        string                   true    "Product Id to be updated"
// @Param       payload       body        dto.UpdateProductInput   true    "Product payload"
// @Success     200           {object}    entity.Product
// @Failure     400           {object}    Error
// @Failure     404           {object}    Error
// @Failure     500           {object}    Error
// @Router      /product/{product_id}     [put]
// @Security    ApiKeyAuth
func (handler *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Println("request does not have an id parameter.")
		w.WriteHeader(http.StatusBadRequest)
		errorObj := Error{Message: "you must provide an id to the url"}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	productID, err := pkgEntity.ParseID(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	var updateProductInput dto.UpdateProductInput
	if err := json.NewDecoder(r.Body).Decode(&updateProductInput); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	product, err := handler.ProductDB.FindById(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	product.ID = productID
	product.Name = updateProductInput.Name
	product.Price = updateProductInput.Price
	if err = handler.ProductDB.Update(product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
	}
}

// Delete product godoc
// @Summary     Delete products
// @Description Delete products by id
// @Tags        products
// @Accept      json
// @Produce     json
// @Param       product_id    path       string   true    "Product Id to be deleted"
// @Success     200
// @Failure     400           {object}    Error
// @Failure     500           {object}    Error
// @Router      /product/{product_id}   [delete]
// @Security    ApiKeyAuth
func (handler *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Println("request does not have an id parameter.")
		w.WriteHeader(http.StatusBadRequest)
		errorObj := Error{Message: "you must provide an id to the url"}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	if err := handler.ProductDB.Delete(id); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// List products godoc
// @Summary     List products
// @Description List all products with support for pagination
// @Tags        products
// @Accept      json
// @Produce     json
// @Param       page          query       string   false    "Page number"
// @Param       limit         query       string   false    "Products per page"
// @Success     200           {array}    entity.Product
// @Failure     500           {object}    Error
// @Router      /product   [get]
// @Security    ApiKeyAuth
func (handler *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 0
	}
	sort := r.URL.Query().Get("sort")

	product, err := handler.ProductDB.FindAll(page, limit, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
	}
}
