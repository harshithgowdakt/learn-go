package main

import (
	"fmt"
	"time"
)

// Models
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Service layer
type UserService struct {
	db     Database
	cache  Cache
	logger Logger
}

func NewUserService(db Database, cache Cache, logger Logger) *UserService {
	return &UserService{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

func (s *UserService) GetUser(id string) (*User, error) {
	s.logger.Info(fmt.Sprintf("Getting user: %s", id))

	// Try cache first
	if cached, found := s.cache.Get("user:" + id); found {
		s.logger.Info("User found in cache")
		return cached.(*User), nil
	}

	// Get from database
	user, err := s.db.GetUser(id)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to get user: %v", err))
		return nil, err
	}

	// Cache the result
	s.cache.Set("user:"+id, user, 5*time.Minute)

	return user, nil
}

func (s *UserService) CreateUser(user *User) error {
	s.logger.Info(fmt.Sprintf("Creating user: %s", user.Name))
	return s.db.CreateUser(user)
}
