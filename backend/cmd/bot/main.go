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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
	
	loggingService := logging.NewService(db.GetClient())
	userService := users.NewService(db.GetClient(), loggingService)
	contentService := content.NewService(
		db.GetClient(), 
		loggingService,
		cfg.ContentBucketName,
	)
	twilioService := twilio.NewService(
		userService,
		contentService,
		cfg.TwilioAccountSid,
		cfg.TwilioAuthToken,
		cfg.TwilioPhoneNumber,
		cfg.TwilioContentSid,
		cfg.TwilioMessagingServiceSid,
	)
	schedulerService := scheduler.NewService(contentService, twilioService, loggingService)
	
	mux := http.NewServeMux()
	
	mux.HandleFunc("/content", contentService.HandleCreateContent)
	mux.HandleFunc("/twilio/webhook", twilioService.HandleWebhook)
	
	schedulerService.Start()
	
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", corsMiddleware(mux)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
