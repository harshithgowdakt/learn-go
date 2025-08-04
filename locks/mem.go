package main

import (
    "encoding/json"
    "net/http"
    "sync"
    "time"
)

type UserCache struct {
    mu    sync.RWMutex
    users map[string]User
    stats CacheStats
}

type User struct {
    ID       string    `json:"id"`
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    LastSeen time.Time `json:"last_seen"`
}

type CacheStats struct {
    Hits   int64 `json:"hits"`
    Misses int64 `json:"misses"`
}

func (c *UserCache) GetUser(id string) (User, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    user, exists := c.users[id]
    if exists {
        c.stats.Hits++
    } else {
        c.stats.Misses++
    }
    return user, exists
}

func (c *UserCache) SetUser(user User) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.users[user.ID] = user
}

func (c *UserCache) GetStats() CacheStats {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    return c.stats
}

// HTTP handlers
func (c *UserCache) userHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    
    user, exists := c.GetUser(userID)
    if !exists {
        // Simulate database lookup
        user = User{
            ID:       userID,
            Name:     "User " + userID,
            Email:    userID + "@example.com",
            LastSeen: time.Now(),
        }
        c.SetUser(user)
    }
    
    json.NewEncoder(w).Encode(user)
}

func main() {
    cache := &UserCache{
        users: make(map[string]User),
    }
    
    http.HandleFunc("/user", cache.userHandler)
    http.ListenAndServe(":8080", nil)
}