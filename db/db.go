package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Connection struct {
	ID       int
	InUse    bool
	LastUsed time.Time
}

type ConnectionPool struct {
	connections chan *Connection
	maxConns    int
	timeout     time.Duration
	mu          sync.Mutex
	created     int
}

func NewConnectionPool(maxConns int, timeout time.Duration) *ConnectionPool {
	return &ConnectionPool{
		connections: make(chan *Connection, maxConns),
		maxConns:    maxConns,
		timeout:     timeout,
	}
}

func (cp *ConnectionPool) GetConnection() (*Connection, error) {
	select {
	case conn := <-cp.connections:
		conn.InUse = true
		conn.LastUsed = time.Now()
		fmt.Printf("♻️ Reused connection %d\n", conn.ID)
		return conn, nil

	default:
		// No available connections, try to create new one
		cp.mu.Lock()
		if cp.created < cp.maxConns {
			cp.created++
			connID := cp.created
			cp.mu.Unlock()

			conn := &Connection{
				ID:       connID,
				InUse:    true,
				LastUsed: time.Now(),
			}
			fmt.Printf("🆕 Created new connection %d\n", conn.ID)
			return conn, nil
		}
		cp.mu.Unlock()

		// Pool is full, wait with timeout
		select {
		case conn := <-cp.connections:
			conn.InUse = true
			conn.LastUsed = time.Now()
			fmt.Printf("⏰ Got connection %d after waiting\n", conn.ID)
			return conn, nil

		case <-time.After(cp.timeout):
			return nil, errors.New("timeout: no connection available")
		}
	}
}

func (cp *ConnectionPool) ReleaseConnection(conn *Connection) {
	conn.InUse = false

	select {
	case cp.connections <- conn:
		fmt.Printf("🔄 Returned connection %d to pool\n", conn.ID)
	default:
		fmt.Printf("⚠️ Pool full, discarding connection %d\n", conn.ID)
	}
}

// Simulate database operation
func performDatabaseOperation(connPool *ConnectionPool, operationID int, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := connPool.GetConnection()
	if err != nil {
		fmt.Printf("❌ Operation %d failed: %v\n", operationID, err)
		return
	}

	// Simulate work
	workTime := time.Duration(operationID%3+1) * time.Second
	fmt.Printf("🔧 Operation %d using connection %d (will take %v)\n",
		operationID, conn.ID, workTime)
	time.Sleep(workTime)

	connPool.ReleaseConnection(conn)
	fmt.Printf("✅ Operation %d completed\n", operationID)
}

func main() {
    // Create pool with max 2 connections and 3 second timeout
    pool := NewConnectionPool(2, 3*time.Second)

    var wg sync.WaitGroup

    // Start 5 concurrent operations (more than pool size)
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go performDatabaseOperation(pool, i, &wg)
        time.Sleep(500 * time.Millisecond) // Stagger starts
    }

    wg.Wait()
    fmt.Println("All operations completed!")
}
