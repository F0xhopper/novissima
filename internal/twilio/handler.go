package twilio

import (
	"net/http"
)

func (s *Service) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	
	if !s.validateRequest(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	from := r.FormValue("From") 
	body := r.FormValue("Body") 
	messageType := r.FormValue("MessageType")

	if messageType != "text" {
		s.sendResponse(w, "Please send a text message to interact with the bot.")
		return
	}
	
	response := s.processMessage(from, body)
	s.sendResponse(w, response)
}
