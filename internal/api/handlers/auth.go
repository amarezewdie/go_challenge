package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/izymalhaw/go-crud/yishakterefe/internal/api/dto"
	"github.com/izymalhaw/go-crud/yishakterefe/internal/services/auth"
	"github.com/izymalhaw/go-crud/yishakterefe/internal/util"
)

type AuthHandler struct {
	authService   *auth.Service
	adminEmail    string
	adminPassword string
}

func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		adminEmail:    os.Getenv("ADMIN_EMAIL"),
		adminPassword: os.Getenv("ADMIN_PASSWORD"),
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//  Parse request body
		var req dto.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Validate required fields
		if req.Email == "" || req.Password == "" {
			util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid credentials")
			return
		}

		// Verify credentials
		if req.Email != h.adminEmail || req.Password != h.adminPassword {
			util.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		//  Generate JWT token
		token, err := h.authService.GenerateToken(os.Getenv("ADMIN_EMAIL"))
		if err != nil {
			util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		// Return successful response
		util.WriteSuccessResponse(w, dto.LoginResponse{Token: token}, "Login successful")
	}
}