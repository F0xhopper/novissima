package main

import (
	"log"
	"net/http"
	"novissima/internal/config"
	"novissima/internal/content"
	"novissima/internal/database"
	"novissima/internal/logging"
	"novissima/internal/scheduler"
	"novissima/internal/twilio"
	"novissima/internal/users"
)

func	 main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	
	db := database.NewSupabaseDB(
		cfg.SupabaseURL,
		cfg.SupabaseKey,
		cfg.SupabaseEmail,
		cfg.SupabasePassword,
	)
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}

	
	loggingService := logging.NewService(db.GetClient())
	userService := users.NewService(db.GetClient(), loggingService)
	contentService := content.NewService(
		db.GetClient(), 
		loggingService,
		cfg.ContentBucketName,
	)
	twilioService := twilio.NewService(
		userService,
		cfg.TwilioAccountSid,
		cfg.TwilioAuthToken,
		cfg.TwilioPhoneNumber,
	)
	schedulerService := scheduler.NewService(contentService, twilioService, loggingService)
	

	http.HandleFunc("/content", contentService.HandleCreateContent)
	http.HandleFunc("/twilio/webhook", twilioService.HandleWebhook)
	
	schedulerService.Start()
	
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
