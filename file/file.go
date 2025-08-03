package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type Job struct {
    ID       int
    Filename string
}

type Result struct {
    Job       Job
    ProcessedSize int
    Duration     time.Duration
}

// Simulates file processing
func processFile(job Job) Result {
    // Simulate work
    processingTime := time.Duration(rand.Intn(3)+1) * time.Second
    time.Sleep(processingTime)
    
    return Result{
        Job:           job,
        ProcessedSize: rand.Intn(1000) + 100,
        Duration:     processingTime,
    }
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for job := range jobs {
        fmt.Printf("Worker %d processing %s\n", id, job.Filename)
        result := processFile(job)
        results <- result
    }
}

func main() {
    files := []string{
        "data1.csv", "data2.csv", "data3.csv", "data4.csv",
        "log1.txt", "log2.txt", "config.json", "backup.sql",
    }
    
    jobs := make(chan Job, len(files))
    results := make(chan Result, len(files))
    
    const numWorkers = 3
    var wg sync.WaitGroup
    
    // Start workers
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go worker(i, jobs, results, &wg)
    }
    
    // Send jobs
    for i, filename := range files {
        jobs <- Job{ID: i + 1, Filename: filename}
    }
    close(jobs)
    
    // Close results when all workers finish
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results
    totalSize := 0
    for result := range results {
        fmt.Printf("✅ Processed %s: %d bytes in %v\n", 
            result.Job.Filename, result.ProcessedSize, result.Duration)
        totalSize += result.ProcessedSize
    }
    
    fmt.Printf("\nTotal processed: %d bytes\n", totalSize)
}