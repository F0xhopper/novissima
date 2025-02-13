package subscriber

import (
	"encoding/json"
	"net/http"
)

func (s *Service) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.AddSubscriber(req.Email); err != nil {
		http.Error(w, "Failed to add subscriber", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
} 