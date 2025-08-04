package main

import (
    "fmt"
    "sync"
    "time"
)

type Event struct {
    Type      string
    Data      interface{}
    Timestamp time.Time
}

type EventBus struct {
    subscribers map[string][]chan Event
    register    chan subscription
    unregister  chan subscription
    events      chan Event
    quit        chan bool
}

type subscription struct {
    eventType string
    channel   chan Event
    remove    bool
}

func NewEventBus() *EventBus {
    eb := &EventBus{
        subscribers: make(map[string][]chan Event),
        register:    make(chan subscription),
        unregister:  make(chan subscription),
        events:      make(chan Event, 100),
        quit:        make(chan bool),
    }
    
    go eb.run()
    
    return eb
}

func (eb *EventBus) run() {
    for {
        select {
        case sub := <-eb.register:
            if !sub.remove {
                eb.subscribers[sub.eventType] = append(eb.subscribers[sub.eventType], sub.channel)
            }
            
        case unsub := <-eb.unregister:
            if subs, ok := eb.subscribers[unsub.eventType]; ok {
                for i, ch := range subs {
                    if ch == unsub.channel {
                        eb.subscribers[unsub.eventType] = append(subs[:i], subs[i+1:]...)
                        close(ch)
                        break
                    }
                }
            }
            
        case event := <-eb.events:
            if subs, ok := eb.subscribers[event.Type]; ok {
                for _, sub := range subs {
                    select {
                    case sub <- event:
                        // Event sent
                    default:
                        // Subscriber not ready, skip
                    }
                }
            }
            
        case <-eb.quit:
            return
        }
    }
}

func (eb *EventBus) Subscribe(eventType string) <-chan Event {
    ch := make(chan Event, 10)
    eb.register <- subscription{
        eventType: eventType,
        channel:   ch,
    }
    return ch
}

func (eb *EventBus) Unsubscribe(eventType string, ch <-chan Event) {
    eb.unregister <- subscription{
        eventType: eventType,
        channel:   ch.(chan Event),
        remove:    true,
    }
}

func (eb *EventBus) Publish(event Event) {
    event.Timestamp = time.Now()
    eb.events <- event
}

func (eb *EventBus) Close() {
    close(eb.quit)
}

// Usage - Order processing system
type OrderProcessor struct {
    eventBus *EventBus
}

func NewOrderProcessor(eventBus *EventBus) *OrderProcessor {
    op := &OrderProcessor{eventBus: eventBus}
    
    // Subscribe to order events
    orderEvents := eventBus.Subscribe("order")
    go op.handleOrderEvents(orderEvents)
    
    return op
}

func (op *OrderProcessor) handleOrderEvents(events <-chan Event) {
    for event := range events {
        switch event.Type {
        case "order":
            fmt.Printf("Processing order: %v\n", event.Data)
            
            // Simulate processing
            time.Sleep(100 * time.Millisecond)
            
            // Publish completion event
            op.eventBus.Publish(Event{
                Type: "order_completed",
                Data: event.Data,
            })
        }
    }
}

func main() {
    eventBus := NewEventBus()
    processor := NewOrderProcessor(eventBus)
    
    // Subscribe to completion events
    completions := eventBus.Subscribe("order_completed")
    go func() {
        for event := range completions {
            fmt.Printf("Order completed: %v at %v\n", event.Data, event.Timestamp)
        }
    }()
    
    // Simulate orders
    for i := 0; i < 5; i++ {
        eventBus.Publish(Event{
            Type: "order",
            Data: fmt.Sprintf("order-%d", i),
        })
    }
    
    time.Sleep(2 * time.Second)
    eventBus.Close()
}