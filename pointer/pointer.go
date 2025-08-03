package pointer

import "fmt"

type Person struct {
    Name string
    Age  int
}

// Struct pointer functions
func updatePerson(p *Person) {
    p.Name = "Updated"  // Go auto-dereferences
    p.Age = 99
}

func createPerson(name string) *Person {
    return &Person{Name: name, Age: 0}
}

// Slice pointer functions
func addToSlice(s *[]int, value int) {
    *s = append(*s, value)
}

func clearSlice(s *[]int) {
    *s = (*s)[:0]  // Keep capacity, set length to 0
}

// Map pointer functions
func addToMap(m *map[string]int, key string, value int) {
    (*m)[key] = value
}

func clearMap(m *map[string]int) {
    *m = make(map[string]int)
}

func TestPointers() {
    // Struct pointers
    person := Person{Name: "John", Age: 25}
    fmt.Println("Before:", person)
    updatePerson(&person)
    fmt.Println("After:", person)
    
    newPerson := createPerson("Alice")
    fmt.Println("New person:", *newPerson)
    
    // Slice pointers
    numbers := []int{1, 2, 3}
    fmt.Println("Before:", numbers)
    addToSlice(&numbers, 4)
    fmt.Println("After:", numbers)
    
    // Map pointers
    data := map[string]int{"a": 1, "b": 2}
    fmt.Println("Before:", data)
    addToMap(&data, "c", 3)
    fmt.Println("After:", data)
}