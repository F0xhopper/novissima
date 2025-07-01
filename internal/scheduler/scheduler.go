package scheduler

import (
	"log"
	"novissima/internal/content"

	"github.com/robfig/cron/v3"
)

type Service struct {
	contentService *content.Service
}	

func NewService(contentService *content.Service) *Service {
	return &Service{
		contentService: contentService,
	}
}

func (s *Service) Start() {
	c := cron.New()
	c.AddFunc("0 8 * * *", func() {
		if err := s.contentService.SendDailyContent(); err != nil {
			log.Printf("Error sending content: %v", err)
		}
	})
	c.Start()
}
