package content

import (
	"log"
	"novissima/internal/subscriber"
)

type Service struct {
	subscriberService *subscriber.Service
}

func NewService(subscriberService *subscriber.Service) *Service {
	return &Service{
		subscriberService: subscriberService,
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
			log.Printf("Error sending meditation to %s: %v", sub.Email, err)
		}
	}
	
	return nil
}

func (s *Service) sendMeditation(subscriber subscriber.Subscriber) error {
	// Implement your meditation sending logic here
	// This could involve:
	// 1. Selecting a meditation from a database
	// 2. Formatting the meditation message
	// 3. Sending via email/SMS/etc.
	return nil
} 