package main

import (
    "encoding/json"
    "net/http"
    "time"
)

type CacheRequest struct {
    Key      string
    Value    interface{}
    Response chan CacheResponse
}

type CacheResponse struct {
    Value  interface{}
    Exists bool
}

type CacheServer struct {
    requests chan CacheRequest
    cache    map[string]interface{}
}

func NewCacheServer() *CacheServer {
    cs := &CacheServer{
        requests: make(chan CacheRequest),
        cache:    make(map[string]interface{}),
    }
    
    // Single goroutine manages all cache operations
    go cs.run()
    
    return cs
}

func (cs *CacheServer) run() {
    for req := range cs.requests {
        if req.Value != nil {
            // Set operation
            cs.cache[req.Key] = req.Value
            req.Response <- CacheResponse{Value: req.Value, Exists: true}
        } else {
            // Get operation
            value, exists := cs.cache[req.Key]
            req.Response <- CacheResponse{Value: value, Exists: exists}
        }
    }
}

func (cs *CacheServer) Get(key string) (interface{}, bool) {
    response := make(chan CacheResponse)
    cs.requests <- CacheRequest{
        Key:      key,
        Response: response,
    }
    
    resp := <-response
    return resp.Value, resp.Exists
}

func (cs *CacheServer) Set(key string, value interface{}) {
    response := make(chan CacheResponse)
    cs.requests <- CacheRequest{
        Key:      key,
        Value:    value,
        Response: response,
    }
    
    <-response // Wait for confirmation
}

func (cs *CacheServer) userHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    
    user, exists := cs.Get(userID)
    if !exists {
        // Simulate database lookup
        user = map[string]interface{}{
            "id":    userID,
            "name":  "User " + userID,
            "email": userID + "@example.com",
        }
        cs.Set(userID, user)
    }
    
    json.NewEncoder(w).Encode(user)
}