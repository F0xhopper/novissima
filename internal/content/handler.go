package content

import (
	"encoding/json"
	"net/http"
)

func (s *Service) HandleCreateContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ContentText string `json:"contentText"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	content, err := s.AddContent(req.ContentText)
	if err != nil {
		http.Error(w, "Failed to add content", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Successfully added content: " + content.Content)
} 