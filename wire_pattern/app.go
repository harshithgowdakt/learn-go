package main

import (
	"fmt"
	"net/http"
)

// Application
type App struct {
	server      *http.Server
	userHandler *UserHandler
	logger      Logger
}

func NewApp(userHandler *UserHandler, logger Logger, port int) *App {
	mux := http.NewServeMux()
	mux.HandleFunc("/user", userHandler.GetUser)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return &App{
		server:      server,
		userHandler: userHandler,
		logger:      logger,
	}
}

func (a *App) Start() error {
	a.logger.Info("Starting server...")
	return a.server.ListenAndServe()
}

// Manual Wire Function - This is the "Wire Pattern"
func InitializeApp(cfg Config) (*App, error) {
	// Initialize database
	db, err := NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize logger
	logger := NewLogger(cfg.LogLevel)

	// Initialize cache
	cache := NewRedisCache()

	// Initialize service with dependencies
	userService := NewUserService(db, cache, logger)

	// Initialize handler with service
	userHandler := NewUserHandler(userService)

	// Initialize app with all components
	app := NewApp(userHandler, logger, cfg.Port)

	return app, nil
}
