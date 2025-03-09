package content

import (
	"fmt"
	"log"
	"time"

	"github.com/supabase-community/supabase-go"
)

type Content struct {
	Content string
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

func (s *Service) AddContent(contentText string) (Content, error) {
	content := Content{
		Content: contentText,
		CreatedAt: time.Now().UTC(),
	}

	data, _, err := s.client.From("content").Insert(content, true, "", "", "").Execute()
	if err != nil {
		log.Printf("Supabase error: %v", err)
		if data != nil {
			log.Printf("Supabase response: %s", string(data))
		}
		return Content{}, fmt.Errorf("failed to add content: %w", err)
	}

	log.Printf("Successfully added content: %s", data)

	return content, nil
}

func (s *Service) GetDailyMeditation() error {	
	return nil
}
