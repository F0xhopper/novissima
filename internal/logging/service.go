package logging

import (
	"encoding/json"
	"time"

	"github.com/supabase-community/supabase-go"
)

type LogEntry struct {
	ID        int64    
	EventType string    
	Message   string    
	Data      string    
	CreatedAt time.Time 
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

func (s *Service) LogContentSent(contentID int64, message string) error {
	return s.LogEvent("content_sent", message, map[string]interface{}{
		"content_id": contentID,
	})
}

func (s *Service) LogUserCreated(userID int64, username string) error {
	return s.LogEvent("user_created", "New user created", map[string]interface{}{
		"user_id":  userID,
		"username": username,
	})
}

func (s *Service) LogUserDeactivated(userID int64) error {
	return s.LogEvent("user_deactivated", "User deactivated", map[string]interface{}{
		"user_id": userID,
	})
}
func (s *Service) LogUserActivated(userID int64) error {
	return s.LogEvent("user_activated", "User activated", map[string]interface{}{
		"user_id": userID,
	})
}
func (s *Service) LogUserLanguageChanged(userID int64, language string) error {
	return s.LogEvent("user_language_changed", "User language changed", map[string]interface{}{
		"user_id":  userID,
		"language": language,
	})
}	

