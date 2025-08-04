package main

import (
    "database/sql"
    "errors"
    "sync"
    "time"
)

type ConnectionPool struct {
    mu          sync.Mutex
    connections []*sql.DB
    available   []bool
    maxConns    int
    activeConns int
}

func NewConnectionPool(maxConns int, dbURL string) *ConnectionPool {
    pool := &ConnectionPool{
        connections: make([]*sql.DB, maxConns),
        available:   make([]bool, maxConns),
        maxConns:    maxConns,
    }
    
    // Initialize connections
    for i := 0; i < maxConns; i++ {
        db, _ := sql.Open("postgres", dbURL)
        pool.connections[i] = db
        pool.available[i] = true
    }
    
    return pool
}

func (p *ConnectionPool) GetConnection() (*sql.DB, error) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    // Find available connection
    for i := 0; i < p.maxConns; i++ {
        if p.available[i] {
            p.available[i] = false
            p.activeConns++
            return p.connections[i], nil
        }
    }
    
    return nil, errors.New("no connections available")
}

func (p *ConnectionPool) ReleaseConnection(db *sql.DB) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    // Find and release the connection
    for i := 0; i < p.maxConns; i++ {
        if p.connections[i] == db {
            p.available[i] = true
            p.activeConns--
            return
        }
    }
}

func (p *ConnectionPool) Stats() (active, total int) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    return p.activeConns, p.maxConns
}