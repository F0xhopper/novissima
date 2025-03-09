package messaging

import (
	"log"
	"novissima/internal/content"
	"novissima/internal/subscriber"

	"github.com/supabase-community/supabase-go"
)

type Service struct {
	subscriberService *subscriber.Service
	contentService *content.Service
	client *supabase.Client
}

func NewService(subscriberService *subscriber.Service, contentService *content.Service, client *supabase.Client) *Service {
	return &Service{
		subscriberService: subscriberService,
		contentService: contentService,
		client: client,
	}
}

func (s *Service) SendDailyMeditations() error {
	subscribers, err := s.subscriberService.GetAllActiveSubscribers()
	if err != nil {
		return err
	}
	
	for _, sub := range subscribers {
		if err := s.sendMeditation(sub); err != nil {
			// Log error but continue with other subscribers
			log.Printf("Error sending meditation to %d: %v", sub.Phone, err)
		}
	}
	
	return nil
}

func (s *Service) sendMeditation(sub subscriber.Subscriber) error {
	// meditation, err := s.contentService.GetDailyMeditation()
	// if err != nil {
	// 	return err
	// }		
	
	return nil
}	