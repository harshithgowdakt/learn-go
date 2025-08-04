package main

import (
    "fmt"
    "strings"
    "time"
)

type LogEntry struct {
    Timestamp time.Time
    Level     string
    Message   string
    Source    string
}

type ProcessedLog struct {
    LogEntry
    WordCount int
    Severity  int
}

type LogProcessor struct {
    input    chan LogEntry
    filtered chan LogEntry
    parsed   chan ProcessedLog
    output   chan ProcessedLog
}

func NewLogProcessor() *LogProcessor {
    lp := &LogProcessor{
        input:    make(chan LogEntry, 100),
        filtered: make(chan LogEntry, 100),
        parsed:   make(chan ProcessedLog, 100),
        output:   make(chan ProcessedLog, 100),
    }
    
    // Start pipeline stages
    go lp.filterStage()
    go lp.parseStage()
    go lp.enrichStage()
    go lp.outputStage()
    
    return lp
}

func (lp *LogProcessor) filterStage() {
    for log := range lp.input {
        // Filter out debug logs
        if log.Level != "DEBUG" {
            lp.filtered <- log
        }
    }
    close(lp.filtered)
}

func (lp *LogProcessor) parseStage() {
    for log := range lp.filtered {
        processed := ProcessedLog{
            LogEntry:  log,
            WordCount: len(strings.Fields(log.Message)),
        }
        
        // Add severity scoring
        switch log.Level {
        case "ERROR":
            processed.Severity = 3
        case "WARN":
            processed.Severity = 2
        case "INFO":
            processed.Severity = 1
        default:
            processed.Severity = 0
        }
        
        lp.parsed <- processed
    }
    close(lp.parsed)
}

func (lp *LogProcessor) enrichStage() {
    for log := range lp.parsed {
        // Add additional enrichment
        if strings.Contains(log.Message, "error") {
            log.Severity += 1
        }
        
        lp.output <- log
    }
    close(lp.output)
}

func (lp *LogProcessor) outputStage() {
    for log := range lp.output {
        fmt.Printf("[%s] %s: %s (words: %d, severity: %d)\n",
            log.Level, log.Source, log.Message, log.WordCount, log.Severity)
    }
}

func (lp *LogProcessor) Process(log LogEntry) {
    lp.input <- log
}

func (lp *LogProcessor) Close() {
    close(lp.input)
}

// Usage
func main() {
    processor := NewLogProcessor()
    
    // Simulate log entries
    logs := []LogEntry{
        {time.Now(), "INFO", "User login successful", "auth"},
        {time.Now(), "ERROR", "Database connection failed", "db"},
        {time.Now(), "DEBUG", "Processing request", "api"},
        {time.Now(), "WARN", "High memory usage detected", "system"},
        {time.Now(), "ERROR", "Payment processing error occurred", "payment"},
    }
    
    for _, log := range logs {
        processor.Process(log)
    }
    
    time.Sleep(1 * time.Second)
    processor.Close()
    time.Sleep(1 * time.Second) // Wait for processing
}