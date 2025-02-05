package http

import (
	"encoding/json"
	"net/http"
)

type sendMessageRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (c *MainController) SendEmailMessage(w http.ResponseWriter, r *http.Request) {
	var body sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if body.From == "" || body.To == "" || body.Subject == "" || body.Message == "" {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	err := c.mailSender.SendEmail(body.From, body.To, body.Subject, body.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
