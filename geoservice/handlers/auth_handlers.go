package handlers

import (
	"encoding/json"
	"geoservice/models"
	"geoservice/storage"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

var userStore = storage.NewUserStorage()

// @Summary Register User
// @Description Регистрация нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Данные для регистрации"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} models.AuthResponse
// @Router /register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	// Сохранение пользователя
	if err := userStore.AddUser(user); err != nil {
		json.NewEncoder(w).Encode(models.AuthResponse{Error: "Пользователь существует"})
		return
	}

	json.NewEncoder(w).Encode(models.AuthResponse{})
}

// @Summary Login User
// @Description Аутентификация и получение JWT токена
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Данные для входа"
// @Success 200 {object} models.AuthResponse
// @Failure 401 {object} models.AuthResponse
// @Router /login [post]
func LoginHandler(tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		// Поиск пользователя
		user, exists := userStore.GetUser(req.Username)
		if !exists {
			resp := models.AuthResponse{Error: "Неверное имя пользователя или пароль"}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Проверка пароля
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			resp := models.AuthResponse{Error: "Неверное имя пользователя или пароль"}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Генерация JWT токена
		claims := map[string]interface{}{
			"user_id": user.Username,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}

		_, tokenString, err := tokenAuth.Encode(claims)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(models.AuthResponse{
			Token: tokenString,
		})
	}
}
