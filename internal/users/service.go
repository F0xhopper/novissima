package users

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

type User struct {
	ID        int64  
	Phone     int       
	Active    bool      
	Language  string    
	UpdatedAt time.Time 
	CreatedAt time.Time 
}

func NewService(client *supabase.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) AddUser(phone int64) (User, error) {
	user := User{
		Phone:     int(phone),
		Active:    true,
		Language:  "en",
		CreatedAt: time.Now().UTC(),
	}

	data, _, err := s.client.From("users").Insert(user, true, "", "", "").Execute()
	if err != nil {
		log.Printf("Supabase error: %v", err)
		if data != nil {
			log.Printf("Supabase response: %s", string(data))
		}
		return User{}, fmt.Errorf("failed to add user: %w", err)
	}
	
	log.Printf("Successfully added user with phone: %d", phone)
	return user, nil
}

func (s *Service) GetAllActiveUsers() ([]User, error) {
	var users []User
	
	data, _, err := s.client.From("users").
		Select("*", "", false).
		Eq("active", "true").
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}

	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to parse users data: %w", err)
	}

	log.Printf("Found %d active users:", len(users))
	for _, user := range users {
		log.Printf("  - Phone: %d, Created: %s", user.Phone, user.CreatedAt.Format(time.RFC3339))
	}
	
	return users, nil
}
