package main

import (
	"log"

	_ "github.com/lib/pq"
)

// Configuration
type Config struct {
	DatabaseURL string
	Port        int
	LogLevel    string
}

func main() {
	config := Config{
		DatabaseURL: "postgres://localhost/myapp",
		Port:        8080,
		LogLevel:    "info",
	}

	app, err := InitializeApp(config)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start app: %v", err)
	}
}
