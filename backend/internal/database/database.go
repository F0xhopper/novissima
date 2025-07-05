package database

import (
	"log"

	"github.com/supabase-community/supabase-go"
)

type Database interface {
	Connect() error
	GetClient() *supabase.Client
}

type SupabaseDB struct {
	url     string
	key     string
	email   string
	password string
	client  *supabase.Client
}

func NewSupabaseDB(url, key, email, password string) *SupabaseDB {
	return &SupabaseDB{
		url: url,
		key: key,
		email: email,
		password: password,
	}
}

func (db *SupabaseDB) Connect() error {
	client, err := supabase.NewClient(db.url, db.key, nil)
	if err != nil {
		return err
	}

	if db.email == "" || db.password == "" {
		log.Fatal("Missing authentication credentials in .env file")
	}

	_, err = client.Auth.SignInWithEmailPassword(db.email, db.password)
	if err != nil {
		log.Fatalf("Error signing in: %v", err)
	}

	db.client = client
	return nil
}

func (db *SupabaseDB) GetClient() *supabase.Client {
	return db.client
} 