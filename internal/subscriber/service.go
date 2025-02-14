package subscriber

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/supabase-community/supabase-go"
)

type Service struct {
	client *supabase.Client
}

type Subscriber struct {
	Phone     int
	Active    bool
	CreatedAt time.Time
}

func NewService(client *supabase.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) AddSubscriber(phone int64) error {
	subscriber := Subscriber{
		Phone:     int(phone),
		Active:    true,
		CreatedAt: time.Now().UTC(), // Explicitly use UTC
	}
	
	// Log the subscriber data being sent
	data, _, err := s.client.From("subscribers").Insert(subscriber, true, "", "", "").Execute()
	if err != nil {
		// Log the raw error and data
		log.Printf("Supabase error: %v", err)
		if data != nil {
			log.Printf("Supabase response: %s", string(data))
		}
		return fmt.Errorf("failed to add subscriber: %w", err)
	}
	
	log.Printf("Successfully added subscriber with phone: %d", phone)
	return nil
}

func (s *Service) GetAllActiveSubscribers() ([]Subscriber, error) {
	var subscribers []Subscriber
	
	data, _, err := s.client.From("subscribers").
		Select("*", "", false).
		Eq("active", "true").
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get active subscribers: %w", err)
	}

	if err := json.Unmarshal(data, &subscribers); err != nil {
		return nil, fmt.Errorf("failed to parse subscribers data: %w", err)
	}

	log.Printf("Found %d active subscribers:", len(subscribers))
	for _, sub := range subscribers {
		log.Printf("  - Phone: %d, Created: %s", sub.Phone, sub.CreatedAt.Format(time.RFC3339))
	}
	
	return subscribers, nil
}
