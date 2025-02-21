package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"` // Хранится в виде хэша
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	Error string `json:"error,omitempty"`
}
