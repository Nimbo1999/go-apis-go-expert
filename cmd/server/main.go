package main

import (
	"fmt"
	"net/http"

	"github.com/Nimbo1999/go-apis-go-expert/configs"
	"github.com/Nimbo1999/go-apis-go-expert/internal/entity"
	"github.com/Nimbo1999/go-apis-go-expert/internal/infra/database"
	"github.com/Nimbo1999/go-apis-go-expert/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	config, err := configs.LoadConfig("./")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Post("/product", productHandler.Create)
	r.Get("/product/{id}", productHandler.GetProduct)
	r.Get("/product", productHandler.GetProducts)
	r.Put("/product/{id}", productHandler.UpdateProduct)
	r.Delete("/product/{id}", productHandler.DeleteProduct)
	// http.HandleFunc("/product", productHandler.Create)

	err = http.ListenAndServe(fmt.Sprintf(":%s", config.WebServerPort), r)
	if err != nil {
		panic(err)
	}
}
