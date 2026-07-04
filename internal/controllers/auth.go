package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Arup3201/gotask/internal/models"
	"github.com/Arup3201/gotask/internal/utils"
)

type AuthController struct {
	userService *models.UserService
	jwtService  *utils.JWTService
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var data RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {

		http.Error(w,
			"json parse error",
			http.StatusBadRequest)
		return
	}

	user, err := ac.userService.CreateUser(r.Context(), data.Email, data.Name, data.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RegisterResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserAvatar struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginResponse struct {
	AccessToken string     `json:"access_token"`
	ExpiresAt   time.Time  `json:"expires_at"`
	User        UserAvatar `json:"user"`
}

func NewAuthController(userService *models.UserService, jwtService *utils.JWTService) *AuthController {
	return &AuthController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var data LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {

		http.Error(w,
			"json parse error",
			http.StatusBadRequest)
		return
	}

	user, err := ac.userService.ExchangeUserWithCredentials(r.Context(), data.Email, data.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := ac.jwtService.GenerateToken(user.ID, data.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{
		AccessToken: token,
		ExpiresAt:   time.Now().Add(utils.TOKEN_EXPIRES_IN),
		User: UserAvatar{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}
