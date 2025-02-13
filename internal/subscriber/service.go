package subscriber

import (
	"database/sql"
	"time"
)

type Service struct {
	db *sql.DB
}

type Subscriber struct {
	Email     string
	Active    bool
	CreatedAt time.Time
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) AddSubscriber(email string) error {
	query := `
		INSERT INTO subscribers (email, active, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE 
		SET active = true, created_at = EXCLUDED.created_at
	`
	_, err := s.db.Exec(query, email, true, time.Now())
	return err
}

func (s *Service) GetAllActiveSubscribers() ([]Subscriber, error) {
	query := `
		SELECT email, active, created_at 
		FROM subscribers 
		WHERE active = true
	`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscribers []Subscriber
	for rows.Next() {
		var sub Subscriber
		if err := rows.Scan(&sub.Email, &sub.Active, &sub.CreatedAt); err != nil {
			return nil, err
		}
		subscribers = append(subscribers, sub)
	}
	return subscribers, rows.Err()
} 