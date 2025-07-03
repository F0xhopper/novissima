package users

import (
	"encoding/json"
	"fmt"
	"log"
	"novissima/internal/logging"
	"time"

	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

type Service struct {
	client *supabase.Client
	loggingService *logging.Service
}

type UserCreate struct {
	PhoneNumber string `json:"phone_number"`
	Active      bool   `json:"active"`
	Language    string `json:"language"`
}

type UserUpdate struct {
	Active    *bool      `json:"active,omitempty"`
	Language  *string    `json:"language,omitempty"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type User struct {
	ID          uuid.UUID `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	Active      bool      `json:"active"`
	Language    string    `json:"language"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewService(client *supabase.Client, loggingService *logging.Service) *Service {
	return &Service{
		client: client,
		loggingService: loggingService,
	}
}

func (s *Service) AddUser(phoneNumber string) (User, error) {
	user := UserCreate{
		PhoneNumber:  phoneNumber,
		Active:       true,
		Language:     "en",
	}

	
	data, _, err := s.client.From("users").Insert(user, true, "", "", "").Execute()
	if err != nil {
		log.Printf("Supabase error: %v", err)
		return User{}, fmt.Errorf("failed to add user: %w", err)
	}

	var createdUsers []User
	if err := json.Unmarshal(data, &createdUsers); err != nil {
		return User{}, fmt.Errorf("failed to parse created user: %w", err)
	}
	
	if len(createdUsers) == 0 {
		return User{}, fmt.Errorf("no user was created")
	}
	
	createdUser := createdUsers[0]
	s.loggingService.LogUserCreated(createdUser.ID, createdUser.PhoneNumber)
	log.Printf("Successfully added user with phone: %s", phoneNumber)
	return createdUser, nil
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
		log.Printf("  - Phone: %s, Created: %s", user.PhoneNumber, user.CreatedAt.Format(time.RFC3339))
	}
	
	return users, nil
}

func (s *Service) GetUserByPhoneNumber(phoneNumber string) (User, error) {
	var user User
	
	data, _, err := s.client.From("users").
		Select("*", "", false).
		Eq("phone_number", phoneNumber).
		Execute()
	if err != nil {
		return User{}, fmt.Errorf("failed to get user by phone number: %w", err)
	}

	if err := json.Unmarshal(data, &user); err != nil {
		return User{}, fmt.Errorf("failed to parse user data: %w", err)
	}
	
	return user, nil
}

func (s *Service) UpdateUserStatus(phoneNumber string, status string) error {
	_, _, err := s.client.From("users").
		Update(map[string]interface{}{"active": status}, "", "").
		Eq("phone_number", phoneNumber).
		Execute()
	return err
}		

func (s *Service) UpdateUserLanguage(phoneNumber string, language string) error {
	_, _, err := s.client.From("users").
		Update(map[string]interface{}{"language": language}, "", "").
		Eq("phone_number", phoneNumber).
		Execute()
	return err
}		
