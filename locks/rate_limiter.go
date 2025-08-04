package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "sync"
    "time"
)

type Config struct {
    Database struct {
        Host     string `json:"host"`
        Port     int    `json:"port"`
        Username string `json:"username"`
    } `json:"database"`
    
    API struct {
        RateLimit int    `json:"rate_limit"`
        Timeout   string `json:"timeout"`
    } `json:"api"`
}

type ConfigManager struct {
    mu         sync.RWMutex
    config     Config
    configPath string
    lastMod    time.Time
}

func NewConfigManager(configPath string) *ConfigManager {
    cm := &ConfigManager{
        configPath: configPath,
    }
    
    cm.loadConfig()
    
    // Start background config watcher
    go cm.watchConfig()
    
    return cm
}

func (cm *ConfigManager) loadConfig() {
    data, err := ioutil.ReadFile(cm.configPath)
    if err != nil {
        log.Printf("Error reading config: %v", err)
        return
    }
    
    var newConfig Config
    if err := json.Unmarshal(data, &newConfig); err != nil {
        log.Printf("Error parsing config: %v", err)
        return
    }
    
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    cm.config = newConfig
    log.Println("Configuration reloaded")
}

func (cm *ConfigManager) GetConfig() Config {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    
    return cm.config
}

func (cm *ConfigManager) GetDatabaseConfig() (string, int, string) {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    
    db := cm.config.Database
    return db.Host, db.Port, db.Username
}

func (cm *ConfigManager) watchConfig() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        // Check if config file was modified
        // In real implementation, you'd use file system events
        cm.loadConfig()
    }
}