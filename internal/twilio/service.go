package twilio

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"novissima/internal/content"
	"novissima/internal/users"
	"sort"
	"strings"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Service struct {
	client      *Client
	userService *users.Service
	accountSid  string
	authToken   string
	phoneNumber string
}

func NewService(userService *users.Service, accountSid, authToken, phoneNumber string) *Service {
	return &Service{
		client:      NewClient(accountSid, authToken, phoneNumber),
		userService: userService,
		accountSid:  accountSid,
		authToken:   authToken,
		phoneNumber: phoneNumber,
	}
}

func (s *Service) SendMessageToUser(phoneNumber, message string, mediaUrl string) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo("whatsapp:" + phoneNumber)
	params.SetFrom("whatsapp:" + s.client.phoneNumber)
	params.SetBody(message)
	params.SetMediaUrl([]string{mediaUrl})

	_, err := s.client.twilioClient.Api.CreateMessage(params)
	return err
}

func (s *Service) SendMessageToAllUsers(content *content.Content) error {
	users, err := s.userService.GetAllActiveUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		var messageToSend string
		
		switch user.Language {
		case "en":
			messageToSend = content.TextEnglish
		case "la":
			messageToSend = content.TextLatin
		case "both":
			messageToSend = content.TextLatin + "\n\n" + content.TextEnglish
		default:
			messageToSend = content.TextEnglish
		}
		
		if err := s.SendMessageToUser(user.PhoneNumber, messageToSend, content.ImageURL); err != nil {
			log.Printf("Error sending message to user %s: %v", user.PhoneNumber, err)
			continue
		}
	}
	
	return nil
}

func (s *Service) validateRequest(r *http.Request) bool {
	
	signature := r.Header.Get("X-Twilio-Signature")
	if signature == "" {
		return false
	}

	
	fullURL := "https://" + r.Host + r.URL.String()
	
	
	values := r.URL.Query()
	for key, val := range r.PostForm {
		values[key] = val
	}

	
	var keys []string
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	
	var buf strings.Builder
	buf.WriteString(fullURL)
	for _, key := range keys {
		buf.WriteString(key)
		for _, val := range values[key] {
			buf.WriteString(val)
		}
	}

	
	h := hmac.New(sha1.New, []byte(s.authToken))
	h.Write([]byte(buf.String()))
	expectedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature == expectedSignature
}

func (s *Service) processMessage(from, body string) string {
	
	cleanNumber := strings.TrimPrefix(from, "whatsapp:")
	
	parts := strings.Fields(strings.TrimSpace(body))
	
	command := strings.ToLower(parts[0])
	
	switch command {
	case "start":
		return s.startSubscription(cleanNumber)
	case "stop":
		return s.stopSubscription(cleanNumber)
	case "help":
		return s.getHelp()
	case "lang":
		return s.setLanguage(cleanNumber, parts)
	case "source":
		return s.getSource(cleanNumber)
	default:
		return "Unknown command. Text 'help' to see available commands."
	}
}

func (s *Service) ensureUserExists(cleanNumber string) (*users.User, error) {
	user, err := s.userService.GetUserByPhoneNumber(cleanNumber)
	if err != nil {
		user, err = s.userService.AddUser(cleanNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to register user: %w", err)
		}
		return &user, nil
	}
	
	return &user, nil
}

func (s *Service) startSubscription(cleanNumber string) string {
	
	user, err := s.ensureUserExists(cleanNumber)
	if err != nil {
		return "Sorry, there was an error starting your subscription. Please try again later."
	}
	
	if !user.Active {
		err = s.userService.UpdateUserStatus(cleanNumber, "active")
		if err != nil {
			return "Sorry, there was an error starting your subscription. Please try again later."
		}
	}
	
	return "Your daily content subscription has been started! You'll receive updates daily."
	
}

func (s *Service) getSource(cleanNumber string) string {

	if _, err := s.ensureUserExists(cleanNumber); err != nil {
		return "Sorry, there was an error getting your source. Please try again later."
	}

	return "Your source is https://novissima.com"

}

func (s *Service) stopSubscription(cleanNumber string) string {
	
	user, err := s.ensureUserExists(cleanNumber)
	if err != nil {
		return "Sorry, there was an error starting your subscription. Please try again later."
	}

	if user.Active {
		err = s.userService.UpdateUserStatus(cleanNumber, "inactive")
		if err != nil {
			return "Sorry, there was an error stopping your subscription. Please try again later."
		}
	}

	return "Your daily content subscription has been stopped. Send 'start' to resume."
}

func (s *Service) getHelp() string {
	return `Available commands:
• start - Start receiving daily content (registers you if needed)
• stop - Stop receiving daily content
• help - Show this help message
• lang [en|la|both] - Set language (English/Latin/Both)
• source - Get the source of the image and text

Example: "start" or "lang en"`
}

func (s *Service) setLanguage(cleanNumber string, parts []string) string {
	
	user, err := s.ensureUserExists(cleanNumber)
	if err != nil {
		return "Sorry, there was an error starting your subscription. Please try again later."
	}
	
	
	if len(parts) < 2 {
		return "Please specify a language. Usage: lang [en|la]"
	}

	language := strings.ToLower(parts[1])
	switch language {
	case "en":
		if user.Language == "en" {
			return "Language is already set to English."
		}
		err = s.userService.UpdateUserLanguage(cleanNumber, "en")
		if err != nil {
			return "Sorry, there was an error updating your language. Please try again later."
		}
		return "Language set to English."
	case "la":
		if user.Language == "la" {
			return "Language is already set to Latin."
		}
		err = s.userService.UpdateUserLanguage(cleanNumber, "la")
		if err != nil {
			return "Sorry, there was an error updating your language. Please try again later."
		}
		return "Language set to Latin."
	case "both":
		if user.Language == "both" {
			return "Language is already set to both."
		}
		err = s.userService.UpdateUserLanguage(cleanNumber, "both")
		if err != nil {
			return "Sorry, there was an error updating your language. Please try again later."
		}
		return "Language set to both."
	default:
		return "Invalid language. Please use 'en' for English or 'la' for Latin."
	}
}

func (s *Service) sendResponse(w http.ResponseWriter, message string) {
	
	twimlResponse := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Message>%s</Message>
</Response>`, message)

	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(twimlResponse))
} 