package logging

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

type LogEntry struct {
	ID        uuid.UUID  `json:"id"`
	EventType string     `json:"event_type"`
	Message   string     `json:"message"`
	Data      string     `json:"data"`
	CreatedAt time.Time  `json:"created_at"`
}

type Service struct {
	client *supabase.Client
}

func NewService(client *supabase.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) LogEvent(eventType, message string, data map[string]interface{}) error {
	dataJSON := "{}"
	if data != nil {
		if jsonBytes, err := json.Marshal(data); err == nil {
			dataJSON = string(jsonBytes)
		}
	}

	_, _, err := s.client.From("logs").Insert(LogEntry{
		EventType: eventType,
		Message:   message,
		Data:      dataJSON,
		CreatedAt: time.Now(),
	}, true, "", "", "").Execute()
	return err
}

func (s *Service) LogContentCreated(contentID uuid.UUID, textEnglish string, textLatin string, imageURL string, theme string, imageSource string, textSource string) error {
	return s.LogEvent("content_created", "New content created", map[string]interface{}{
		"content_id": contentID,
		"text_english": textEnglish,
		"text_latin": textLatin,
		"image_url": imageURL,
		"theme": theme,
		"image_source": imageSource,
		"text_source": textSource,
	})
}

func (s *Service) LogContentSent(contentID uuid.UUID) error {
	return s.LogEvent("content_sent", "Content sent", map[string]interface{}{
		"content_id": contentID,
	})
}

func (s *Service) LogUserCreated(userID uuid.UUID, phoneNumber string) error {
	return s.LogEvent("user_created", "New user created", map[string]interface{}{
		"user_id":  userID,
		"phone_number": phoneNumber,
	})
}

func (s *Service) LogUserDeactivated(userID uuid.UUID) error {
	return s.LogEvent("user_deactivated", "User deactivated", map[string]interface{}{
		"user_id": userID,
	})
}

func (s *Service) LogUserActivated(userID uuid.UUID) error {
	return s.LogEvent("user_activated", "User activated", map[string]interface{}{
		"user_id": userID,
	})
}
func (s *Service) LogUserLanguageChanged(userID uuid.UUID, language string) error {
	return s.LogEvent("user_language_changed", "User language changed", map[string]interface{}{
		"user_id":  userID,
		"language": language,
	})
}	

