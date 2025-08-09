package main

import (
	"fmt"
	"sync"
)

type LRUCache struct {
	capacity int
	cache    map[string]*Node
	head     *Node
	tail     *Node
	mu       sync.RWMutex
}

type Node struct {
	key   string
	value interface{}
	prev  *Node
	next  *Node
}

// NewLRUCache creates a new LRU cache with given capacity
func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		panic("capacity must be positive")
	}

	// Create dummy head and tail nodes
	head := &Node{}
	tail := &Node{}

	// Connect head and tail
	head.next = tail
	tail.prev = head

	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*Node),
		head:     head,
		tail:     tail,
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.cache[key]; exists {
		c.moveToHead(node)
		return node.value, true
	}
	return nil, false
}

func (c *LRUCache) Put(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.cache[key]; exists {
		node.value = value
		c.moveToHead(node)
		return
	}

	newNode := &Node{key: key, value: value}
	c.cache[key] = newNode
	c.addToHead(newNode)

	if len(c.cache) > c.capacity {
		tail := c.removeTail()
		delete(c.cache, tail.key)
	}
}

// addToHead adds node right after head
func (c *LRUCache) addToHead(node *Node) {
	node.prev = c.head
	node.next = c.head.next

	c.head.next.prev = node
	c.head.next = node
}

// removeNode removes an existing node from the linked list
func (c *LRUCache) removeNode(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// moveToHead moves existing node to head (mark as recently used)
func (c *LRUCache) moveToHead(node *Node) {
	c.removeNode(node)
	c.addToHead(node)
}

// removeTail removes the last node (least recently used)
func (c *LRUCache) removeTail() *Node {
	lastNode := c.tail.prev
	c.removeNode(lastNode)
	return lastNode
}

// Size returns current number of items in cache
func (c *LRUCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

// Clear removes all items from cache
func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Reset the map
	c.cache = make(map[string]*Node)

	// Reset the linked list
	c.head.next = c.tail
	c.tail.prev = c.head
}

// Keys returns all keys in order from most recently used to least recently used
func (c *LRUCache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.cache))
	current := c.head.next

	for current != c.tail {
		keys = append(keys, current.key)
		current = current.next
	}

	return keys
}

// Delete removes a key from the cache
func (c *LRUCache) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.cache[key]; exists {
		c.removeNode(node)
		delete(c.cache, key)
		return true
	}
	return false
}

// Peek gets a value without marking it as recently used
func (c *LRUCache) Peek(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if node, exists := c.cache[key]; exists {
		return node.value, true
	}
	return nil, false
}

// Example usage and testing
func main() {
	// Create cache with capacity 3
	cache := NewLRUCache(3)

	// Test basic operations
	fmt.Println("=== Basic Operations ===")
	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	fmt.Println("Keys after adding a,b,c:", cache.Keys()) // [c, b, a]

	// Access 'a' to make it most recently used
	val, exists := cache.Get("a")
	fmt.Printf("Get 'a': %v, exists: %v\n", val, exists)
	fmt.Println("Keys after accessing 'a':", cache.Keys()) // [a, c, b]

	// Add 'd' - should evict 'b' (least recently used)
	cache.Put("d", 4)
	fmt.Println("Keys after adding 'd':", cache.Keys()) // [d, a, c]

	// Try to get evicted key 'b'
	val, exists = cache.Get("b")
	fmt.Printf("Get 'b': %v, exists: %v\n", val, exists) // nil, false

	fmt.Println("\n=== Update Existing Key ===")
	cache.Put("a", 100) // Update existing key
	val, _ = cache.Get("a")
	fmt.Printf("Updated value of 'a': %v\n", val) // 100

	fmt.Println("\n=== Cache State ===")
	fmt.Printf("Size: %d\n", cache.Size())
	fmt.Printf("All keys: %v\n", cache.Keys())

	fmt.Println("\n=== Peek vs Get ===")
	cache.Put("x", 99)
	cache.Put("y", 98)
	fmt.Println("Before peek:", cache.Keys())

	// Peek doesn't change order
	cache.Peek("d")
	fmt.Println("After peek 'd':", cache.Keys())

	// Get changes order
	cache.Get("d")
	fmt.Println("After get 'd':", cache.Keys())

	fmt.Println("\n=== Delete Operation ===")
	deleted := cache.Delete("x")
	fmt.Printf("Deleted 'x': %v\n", deleted)
	fmt.Println("Keys after delete:", cache.Keys())

	fmt.Println("\n=== Clear Cache ===")
	cache.Clear()
	fmt.Printf("Size after clear: %d\n", cache.Size())
	fmt.Printf("Keys after clear: %v\n", cache.Keys())
}
