package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Nimbo1999/go-apis-go-expert/internal/dto"
	"github.com/Nimbo1999/go-apis-go-expert/internal/entity"
	"github.com/Nimbo1999/go-apis-go-expert/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	UserDB       database.UserInterface
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

func NewUserHandler(db database.UserInterface, Jwt *jwtauth.JWTAuth, JwtExpiresIn int) *UserHandler {
	return &UserHandler{
		UserDB:       db,
		Jwt:          Jwt,
		JwtExpiresIn: JwtExpiresIn,
	}
}

func (handler *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createUserInput dto.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&createUserInput); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := entity.NewUser(createUserInput.Name, createUserInput.Email, createUserInput.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = handler.UserDB.Create(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var getJwtInput dto.GetJwtInput
	if err := json.NewDecoder(r.Body).Decode(&getJwtInput); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := handler.UserDB.FindByEmail(getJwtInput.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !user.ValidatePassword(getJwtInput.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := generateToken(user, handler)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct {
		AccessToken string `json:"access_token"`
	}{AccessToken: token}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func generateToken(user *entity.User, handler *UserHandler) (string, error) {
	claims := make(map[string]interface{})
	jwtauth.SetExpiryIn(claims, time.Second*time.Duration(handler.JwtExpiresIn))
	jwtauth.SetIssuedNow(claims)
	claims["sub"] = user.ID.String()
	claims["email"] = user.Email
	claims["name"] = user.Name

	_, tokenString, err := handler.Jwt.Encode(claims)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
