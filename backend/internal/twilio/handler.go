package twilio

import (
	"log"
	"net/http"
)

func (s *Service) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Println("Recieved webhook from Twilio:", r.FormValue("From"), r.FormValue("Body"), r.FormValue("MessageType") )
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
