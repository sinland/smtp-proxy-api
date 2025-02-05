package http

import (
	"encoding/json"
	"net/http"
)

type sendTgMessageRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func (c *MainController) SendTelegramMessage(w http.ResponseWriter, r *http.Request) {
	var body sendTgMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if body.To == "" || body.Message == "" {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	err := c.tgService.SendMessage(r.Context(), body.To, body.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
