package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Database layer
type Database interface {
	GetUser(id string) (*User, error)
	CreateUser(user *User) error
}

type PostgresDB struct {
	conn *sql.DB
}

func NewPostgresDB(databaseURL string) (*PostgresDB, error) {
	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{conn: conn}, nil
}

func (db *PostgresDB) GetUser(id string) (*User, error) {
	// Implementation here
	return &User{ID: id, Name: "John Doe"}, nil
}

func (db *PostgresDB) CreateUser(user *User) error {
	// Implementation here
	return nil
}
