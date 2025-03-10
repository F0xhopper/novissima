package subscriber

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *Service) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Phone int `json:"phone"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	subscriber, err := s.AddSubscriber(int64(req.Phone))
	if err != nil {
		http.Error(w, "Failed to add subscriber", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Successfully added subscriber: " + strconv.Itoa(subscriber.Phone))
} 