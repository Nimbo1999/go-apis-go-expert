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
	"github.com/go-chi/jwtauth"
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
	userDB := database.NewUser(db)
	productHandler := handlers.NewProductHandler(productDB)
	userHandler := handlers.NewUserHandler(userDB, config.TokenAuth, config.JWTExpiresIn)

	r := chi.NewRouter()
	// r.Use(LogRequest)
	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(middleware.Recoverer)

	r.Route("/product", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.Create)
		r.Get("/{id}", productHandler.GetProduct)
		r.Get("/", productHandler.GetProducts)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Post("/user", userHandler.Create)
	r.Post("/user/login", userHandler.Login)

	err = http.ListenAndServe(fmt.Sprintf(":%s", config.WebServerPort), r)
	if err != nil {
		panic(err)
	}
}

// Creating a middleware
// func LogRequest(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Printf("%s: %s %s%s %s\n", time.Now().Format(time.RFC3339), r.Method, r.Host, r.URL.String(), r.Proto)
// 		next.ServeHTTP(w, r)
// 	})
// }
