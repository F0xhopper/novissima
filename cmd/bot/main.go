package main

import (
	"log"
	"net/http"
	"novissima/internal/config"
	"novissima/internal/content"
	"novissima/internal/database"
	"novissima/internal/subscriber"

	"github.com/robfig/cron/v3"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	db := database.NewPostgresDB(cfg.DatabaseURL)
	if err := db.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	subscriberService := subscriber.NewService(db.GetDB())
	contentService := content.NewService(subscriberService)
	
	// Set up HTTP handlers
	http.HandleFunc("/subscribers", subscriberService.HandleSubscribe)
	
	// Initialize cron job
	c := cron.New()
	
	// Schedule meditation sending every day at 8:00 AM
	_, err = c.AddFunc("0 8 * * *", func() {
		if err := contentService.SendDailyMeditations(); err != nil {
			log.Printf("Error sending content: %v", err)
		}
	})
	
	if err != nil {
		log.Fatalf("Error setting up cron job: %v", err)
	}
	
	// Start the cron scheduler
	c.Start()
	
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
