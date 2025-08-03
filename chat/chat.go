package main

import (
    "fmt"
    "time"
)

type Message struct {
    User    string
    Content string
    Time    time.Time
}

type ChatRoom struct {
    clients   map[string]chan Message
    join      chan string
    leave     chan string
    messages  chan Message
    quit      chan bool
}

func NewChatRoom() *ChatRoom {
    return &ChatRoom{
        clients:  make(map[string]chan Message),
        join:     make(chan string),
        leave:    make(chan string),
        messages: make(chan Message),
        quit:     make(chan bool),
    }
}

func (cr *ChatRoom) Run() {
    fmt.Println("Chat room started...")
    
    for {
        select {
        case user := <-cr.join:
            cr.clients[user] = make(chan Message, 100)
            fmt.Printf("📥 %s joined the chat\n", user)
            
            // Notify all clients
            joinMsg := Message{
                User:    "System",
                Content: fmt.Sprintf("%s joined the chat", user),
                Time:    time.Now(),
            }
            cr.broadcast(joinMsg)
            
        case user := <-cr.leave:
            if clientChan, exists := cr.clients[user]; exists {
                close(clientChan)
                delete(cr.clients, user)
                fmt.Printf("📤 %s left the chat\n", user)
                
                // Notify all clients
                leaveMsg := Message{
                    User:    "System",
                    Content: fmt.Sprintf("%s left the chat", user),
                    Time:    time.Now(),
                }
                cr.broadcast(leaveMsg)
            }
            
        case msg := <-cr.messages:
            fmt.Printf("💬 [%s] %s: %s\n", 
                msg.Time.Format("15:04:05"), msg.User, msg.Content)
            cr.broadcast(msg)
            
        case <-cr.quit:
            fmt.Println("Shutting down chat room...")
            for user, clientChan := range cr.clients {
                close(clientChan)
                fmt.Printf("Disconnected %s\n", user)
            }
            return
            
        case <-time.After(30 * time.Second):
            fmt.Println("💭 Chat room is quiet...")
        }
    }
}

func (cr *ChatRoom) broadcast(msg Message) {
    for user, clientChan := range cr.clients {
        select {
        case clientChan <- msg:
            // Message sent successfully
        default:
            // Client channel is full, remove them
            fmt.Printf("⚠️ Removing unresponsive client: %s\n", user)
            close(clientChan)
            delete(cr.clients, user)
        }
    }
}

// Simulate client behavior
func simulateClient(name string, chatRoom *ChatRoom) {
    // Join chat
    chatRoom.join <- name
    
    // Send some messages
    messages := []string{
        "Hello everyone!",
        "How's everyone doing?",
        "This is a great chat room!",
    }
    
    for i, content := range messages {
        time.Sleep(time.Duration(i+1) * time.Second)
        chatRoom.messages <- Message{
            User:    name,
            Content: content,
            Time:    time.Now(),
        }
    }
    
    // Leave after some time
    time.Sleep(2 * time.Second)
    chatRoom.leave <- name
}

func main() {
    chatRoom := NewChatRoom()
    
    // Start chat room in background
    go chatRoom.Run()
    
    // Simulate clients joining
    go simulateClient("Alice", chatRoom)
    go simulateClient("Bob", chatRoom)
    go simulateClient("Charlie", chatRoom)
    
    // Let it run for a while
    time.Sleep(15 * time.Second)
    
    // Shutdown
    chatRoom.quit <- true
    time.Sleep(1 * time.Second)
}