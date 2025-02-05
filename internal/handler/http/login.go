package http

import (
	"encoding/json"
	"github.com/sinland/smtp-proxy-api/internal/domain"
	"net/http"
)

type loginRequest struct {
	Username string `json:"username"`
}

type loginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
}

func (c *MainController) Login(w http.ResponseWriter, r *http.Request) {
	var body loginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	jwtToken, err := domain.GenerateToken(body.Username, c.appConfig.Server.JwtSecret)
	if err != nil {
		http.Error(w, "JWT failed", http.StatusInternalServerError)
		return
	}

	res := loginResponse{
		Token:     jwtToken,
		TokenType: "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
