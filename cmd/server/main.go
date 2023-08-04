package main

import (
	"fmt"
	"net/http"

	"github.com/Nimbo1999/go-apis-go-expert/configs"
	"github.com/Nimbo1999/go-apis-go-expert/internal/entity"
	"github.com/Nimbo1999/go-apis-go-expert/internal/infra/database"
	"github.com/Nimbo1999/go-apis-go-expert/internal/infra/webserver/handlers"
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

	http.HandleFunc("/product", productHandler.Create)

	err = http.ListenAndServe(fmt.Sprintf(":%s", config.WebServerPort), nil)
	if err != nil {
		panic(err)
	}
}
