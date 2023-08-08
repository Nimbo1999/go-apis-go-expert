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

type Error struct {
	Message string `json:"message"`
}

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

// Create user godoc
// @Summary     Create user
// @Description Creates a new user
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       request   body      dto.CreateUserInput   true    "User request payload"
// @Success     201
// @Failure     400       {object}  Error
// @Failure     500       {object}  Error
// @Router      /user    [post]
func (handler *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createUserInput dto.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&createUserInput); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	user, err := entity.NewUser(createUserInput.Name, createUserInput.Email, createUserInput.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	if err = handler.UserDB.Create(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login godoc
// @Summary     Get a user JWT
// @Description Authenticate to server and receive a access token as response
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       request       body        dto.GetJWTInput   true    "user credentials"
// @Success     200           {object}    dto.GetJWTOutput
// @Failure     404
// @Failure     500           {object}    Error
// @Router      /user/login   [post]
func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var getJwtInput dto.GetJWTInput
	if err := json.NewDecoder(r.Body).Decode(&getJwtInput); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	user, err := handler.UserDB.FindByEmail(getJwtInput.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	if !user.ValidatePassword(getJwtInput.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		errorObj := Error{Message: "email or password doesn't exists"}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	token, err := generateToken(user, handler)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dto.GetJWTOutput{AccessToken: token}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorObj := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorObj)
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
