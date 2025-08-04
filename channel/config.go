package main

import (
    "encoding/json"
    "log"
    "os"
    "time"
)

type Config struct {
    Database DatabaseConfig `json:"database"`
    API      APIConfig      `json:"api"`
}

type DatabaseConfig struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

type APIConfig struct {
    RateLimit int `json:"rate_limit"`
    Timeout   int `json:"timeout"`
}

type ConfigManager struct {
    configPath    string
    configUpdates chan Config
    subscribers   []chan Config
    current       Config
}

func NewConfigManager(configPath string) *ConfigManager {
    cm := &ConfigManager{
        configPath:    configPath,
        configUpdates: make(chan Config),
        subscribers:   make([]chan Config, 0),
    }
    
    // Load initial config
    cm.loadConfig()
    
    // Start config manager
    go cm.run()
    
    // Start file watcher
    go cm.watchFile()
    
    return cm
}

func (cm *ConfigManager) run() {
    for newConfig := range cm.configUpdates {
        cm.current = newConfig
        log.Println("Configuration updated")
        
        // Notify all subscribers
        for _, subscriber := range cm.subscribers {
            select {
            case subscriber <- newConfig:
                // Sent successfully
            default:
                // Subscriber not ready, skip
                log.Println("Subscriber not ready for config update")
            }
        }
    }
}

func (cm *ConfigManager) loadConfig() {
    data, err := os.ReadFile(cm.configPath)
    if err != nil {
        log.Printf("Error reading config: %v", err)
        return
    }
    
    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        log.Printf("Error parsing config: %v", err)
        return
    }
    
    select {
    case cm.configUpdates <- config:
        // Config update sent
    default:
        // Channel full, skip
    }
}

func (cm *ConfigManager) watchFile() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        cm.loadConfig()
    }
}

func (cm *ConfigManager) Subscribe() <-chan Config {
    subscriber := make(chan Config, 1)
    cm.subscribers = append(cm.subscribers, subscriber)
    
    // Send current config immediately
    subscriber <- cm.current
    
    return subscriber
}

func (cm *ConfigManager) GetCurrent() Config {
    return cm.current
}

// Usage - Database connection manager
type DatabaseManager struct {
    config     Config
    configChan <-chan Config
}

func NewDatabaseManager(cm *ConfigManager) *DatabaseManager {
    dm := &DatabaseManager{
        config:     cm.GetCurrent(),
        configChan: cm.Subscribe(),
    }
    
    go dm.handleConfigUpdates()
    
    return dm
}

func (dm *DatabaseManager) handleConfigUpdates() {
    for newConfig := range dm.configChan {
        if newConfig.Database != dm.config.Database {
            log.Println("Database config changed, reconnecting...")
            dm.config = newConfig
            // Reconnect to database with new config
        }
    }
}