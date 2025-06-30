package main

import (
	"log"
	"net/http"
	"novissima/internal/config"
	"novissima/internal/content"
	"novissima/internal/database"
	"novissima/internal/messaging"
	"novissima/internal/users"

	"github.com/robfig/cron/v3"
)

func main() {

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

	userService := users.NewService(db.GetClient())
	contentService := content.NewService(db.GetClient())
	messagingService := messaging.NewService(userService, contentService, db.GetClient())

	http.HandleFunc("/users", userService.HandleAddUser)
	http.HandleFunc("/content", contentService.HandleCreateContent)
	c := cron.New()
	
	_, err = c.AddFunc("0 8 * * *", func() {
		if err := messagingService.SendDailyContent(); err != nil {
			log.Printf("Error sending content: %v", err)
		}
	})
	
	if err != nil {
		log.Fatalf("Error setting up cron job: %v", err)
	}
	
	c.Start()
	
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
