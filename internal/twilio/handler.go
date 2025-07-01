package twilio

import (
	"fmt"
	"net/http"
	"novissima/internal/users"
	"strings"
)

func (s *Service) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Validate the request is from Twilio
	if !s.validateRequest(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	from := r.FormValue("From") // WhatsApp number
	body := r.FormValue("Body") // Message content
	messageType := r.FormValue("MessageType")

	// Only process text messages
	if messageType != "text" {
		s.sendResponse(w, "Please send a text message to register.")
		return
	}

	// Handle user registration
	response := s.handleRegistration(from, body)
	s.sendResponse(w, response)
}

func (s *Service) validateRequest(r *http.Request) bool {
	signature := r.Header.Get("X-Twilio-Signature")
	url := fmt.Sprintf("https://%s%s", r.Host, r.URL.String())
	
	params := make(map[string]string)
	for key, values := range r.URL.Query() {
		params[key] = values[0]
	}
	
	return s.validator.Validate(url, params, signature)
}

func (s *Service) handleRegistration(from, body string) string {
	// Clean the WhatsApp number (remove whatsapp: prefix)
	cleanNumber := strings.TrimPrefix(from, "whatsapp:")
	
	// Check if user already exists
	user, err := s.userService.GetUserByPhoneNumber(cleanNumber)
	if err == nil && user.PhoneNumber != "" {
		return "You are already registered! Send 'help' for available commands."
	}

	// Parse registration command
	// Example: "register John Doe" or just "register"
	parts := strings.Fields(body)
	if len(parts) == 0 || strings.ToLower(parts[0]) != "register" {
		return "Welcome! To register, send 'register [your name]' or just 'register' to use your phone number as your name."
	}

	var name string
	if len(parts) > 1 {
		name = strings.Join(parts[1:], " ")
	} else {
		name = cleanNumber // Use phone number as name if not provided
	}

	// Create new user
	newUser := &users.User{
		PhoneNumber: cleanNumber,
		Name:        name,
		// Add other required fields
	}

	_, err = s.userService.AddUser(newUser.PhoneNumber)
	if err != nil {
		return "Sorry, there was an error registering you. Please try again later."
	}

	return fmt.Sprintf("Welcome %s! You have been successfully registered. You'll receive daily content updates.", name)
}

func (s *Service) sendResponse(w http.ResponseWriter, message string) {
	// Create TwiML response
	say := &twiml.MessagingMessage{
		Body: message,
	}
	
	response := &twiml.MessagingResponse{
		InnerElements: []twiml.Element{say},
	}
	
	twimlResult, err := twiml.Messaging(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(twimlResult))
} 