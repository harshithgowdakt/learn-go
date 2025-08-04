package main

import (
    "fmt"
    "log"
    "time"
)

type Task struct {
    ID       int
    Data     string
    Result   chan TaskResult
    Priority int
}

type TaskResult struct {
    ID     int
    Output string
    Error  error
}

type WorkerPool struct {
    tasks       chan Task
    results     chan TaskResult
    workerCount int
    quit        chan bool
}

func NewWorkerPool(workerCount int, bufferSize int) *WorkerPool {
    wp := &WorkerPool{
        tasks:       make(chan Task, bufferSize),
        results:     make(chan TaskResult, bufferSize),
        workerCount: workerCount,
        quit:        make(chan bool),
    }
    
    // Start workers
    for i := 0; i < workerCount; i++ {
        go wp.worker(i)
    }
    
    // Start result collector
    go wp.resultCollector()
    
    return wp
}

func (wp *WorkerPool) worker(id int) {
    for {
        select {
        case task := <-wp.tasks:
            log.Printf("Worker %d processing task %d", id, task.ID)
            
            // Simulate work
            time.Sleep(time.Duration(100+task.Priority*50) * time.Millisecond)
            
            result := TaskResult{
                ID:     task.ID,
                Output: fmt.Sprintf("Processed: %s by worker %d", task.Data, id),
            }
            
            // Send result back
            if task.Result != nil {
                task.Result <- result
            } else {
                wp.results <- result
            }
            
        case <-wp.quit:
            log.Printf("Worker %d shutting down", id)
            return
        }
    }
}

func (wp *WorkerPool) resultCollector() {
    for result := range wp.results {
        log.Printf("Task %d completed: %s", result.ID, result.Output)
        // Could save to database, send notification, etc.
    }
}

func (wp *WorkerPool) Submit(task Task) {
    wp.tasks <- task
}

func (wp *WorkerPool) SubmitAndWait(task Task) TaskResult {
    task.Result = make(chan TaskResult)
    wp.tasks <- task
    return <-task.Result
}

func (wp *WorkerPool) Shutdown() {
    close(wp.quit)
    close(wp.tasks)
    close(wp.results)
}

// Usage
func main() {
    pool := NewWorkerPool(3, 10)
    
    // Submit async tasks
    for i := 0; i < 5; i++ {
        pool.Submit(Task{
            ID:       i,
            Data:     fmt.Sprintf("async-task-%d", i),
            Priority: i % 3,
        })
    }
    
    // Submit sync task
    result := pool.SubmitAndWait(Task{
        ID:   100,
        Data: "sync-task",
    })
    fmt.Printf("Sync result: %s\n", result.Output)
    
    time.Sleep(2 * time.Second)
    pool.Shutdown()
}