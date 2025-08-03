package main

import (
    "fmt"
    "net/http"
    "sync"
    "time"
)

type Result struct {
    URL    string
    Status int
    Error  error
}

func checkURL(url string, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    
    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get(url)
    
    result := Result{URL: url}
    if err != nil {
        result.Error = err
    } else {
        result.Status = resp.StatusCode
        resp.Body.Close()
    }
    
    results <- result
}

func main() {
    urls := []string{
        "https://google.com",
        "https://github.com",
        "https://stackoverflow.com",
        "https://invalid-url-12345.com",
        "https://golang.org",
    }
    
    results := make(chan Result, len(urls))
    var wg sync.WaitGroup
    
    // Start goroutines for each URL
    for _, url := range urls {
        wg.Add(1)
        go checkURL(url, results, &wg)
    }
    
    // Close channel when all goroutines finish
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results
    fmt.Println("Checking URLs...")
    for result := range results {
        if result.Error != nil {
            fmt.Printf("❌ %s: ERROR - %v\n", result.URL, result.Error)
        } else {
            fmt.Printf("✅ %s: %d\n", result.URL, result.Status)
        }
    }
}