package messaging

import (
	"log"
	"novissima/internal/content"
	"novissima/internal/users"

	"github.com/supabase-community/supabase-go"
)

type Service struct {
	userService *users.Service
	contentService *content.Service
	client *supabase.Client
}

func NewService(userService *users.Service, contentService *content.Service, client *supabase.Client) *Service {
	return &Service{
		userService: userService,
		contentService: contentService,
		client: client,
	}
}

func (s *Service) SendDailyContent() error {
	users, err := s.userService.GetAllActiveUsers()
	if err != nil {
		return err
	}
	
	for _, user := range users {
		if err := s.sendContent(user); err != nil {
			// Log error but continue with other users
			log.Printf("Error sending meditation to %d: %v", user.Phone, err)
		}
	}
	
	return nil
}

func (s *Service) sendContent(user users.User) error {
	// meditation, err := s.contentService.GetDailyMeditation()
	// if err != nil {
	// 	return err
	// }		
	
	return nil
}	