package scheduler

import (
	"log"
	"novissima/internal/content"
	"novissima/internal/logging"
	"novissima/internal/twilio"

	"github.com/robfig/cron/v3"
)

type Service struct {
	contentService *content.Service
	twilioService *twilio.Service
	loggingService *logging.Service
}	

func NewService(contentService *content.Service, twilioService *twilio.Service, loggingService *logging.Service) *Service {
	return &Service{
		contentService: contentService,
		twilioService: twilioService,
		loggingService: loggingService,
	}
}

func (s *Service) Start() {
	c := cron.New()
	
	c.AddFunc("0 8 * * *", func() {
		log.Println("Starting daily content distribution...")
		
		content, err := s.contentService.GetDailyContent()
		if err != nil {
			log.Printf("Error getting daily content: %v", err)
			return
		}
		
		err = s.twilioService.SendMessageToAllUsers(content)
		if err != nil {
			log.Printf("Failed to send daily content: %v", err)
			return
		}
		s.loggingService.LogContentSent(content.ID)
		log.Println("Daily content sent to all users successfully")
	})
	
	c.Start()
	log.Println("Scheduler started - daily content will be sent at 8 AM")
}
