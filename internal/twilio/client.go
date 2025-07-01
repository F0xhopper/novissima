package twilio

import (
	"novissima/internal/content"
	"novissima/internal/users"

	"github.com/twilio/twilio-go/client"
)

type Service struct {
	userService    *users.Service
	contentService *content.Service
	validator      client.RequestValidator
	phoneNumber    string
}

func NewService(userService *users.Service, contentService *content.Service, authToken, phoneNumber string) *Service {
	return &Service{
		userService:    userService,
		contentService: contentService,
		validator:      client.NewRequestValidator(authToken),
		phoneNumber:    phoneNumber,
	}
}
