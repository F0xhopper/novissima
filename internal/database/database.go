package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Database interface {
	Connect() error
	Close() error
	GetDB() *sql.DB
}

type PostgresDB struct {
	connectionString string
	db              *sql.DB
}

func NewPostgresDB(connectionString string) *PostgresDB {
	return &PostgresDB{
		connectionString: connectionString,
	}
}

func (db *PostgresDB) Connect() error {
	conn, err := sql.Open("postgres", db.connectionString)
	if err != nil {
		return err
	}
	
	if err := conn.Ping(); err != nil {
		return err
	}
	
	db.db = conn
	return nil
}

func (db *PostgresDB) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

func (db *PostgresDB) GetDB() *sql.DB {
	return db.db
} 