package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"task-manager/internals/auth"
	"task-manager/internals/models"
	"task-manager/internals/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Users     *repositories.UserRepository
	JWTSecret string
}

func NewAuthHandler(users *repositories.UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{Users: users, JWTSecret: jwtSecret}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if len(req.Username) < 3 || len(req.Password) < 6 {
		http.Error(w, "username must be >= 3 chars and password must be >= 6 chars", http.StatusBadRequest)
		return
	}

	existing, err := h.Users.GetByUsername(r.Context(), req.Username)
	if err != nil {
		http.Error(w, "failed to check user", http.StatusInternalServerError)
		return
	}
	if existing != nil {
		http.Error(w, "username already exists", http.StatusConflict)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to secure password", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: string(hash),
	}
	if err := h.Users.CreateUser(r.Context(), user); err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Username, h.JWTSecret, 24*time.Hour)
	if err != nil {
		http.Error(w, "failed to create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(models.AuthResponse{Token: token, Username: user.Username})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.Users.GetByUsername(r.Context(), req.Username)
	if err != nil {
		http.Error(w, "failed to authenticate", http.StatusInternalServerError)
		return
	}
	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Username, h.JWTSecret, 24*time.Hour)
	if err != nil {
		http.Error(w, "failed to create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(models.AuthResponse{Token: token, Username: user.Username})
}
