package audit

import (
    "fmt"
    "sync"
)

type Decision struct {
    SubjectID string
    Action    string
    Resource  string
    Allowed   bool
    Actor     string
}

type Auditor interface {
    Record(Decision)
}

type Logger struct {
    mu sync.Mutex
}

func NewLogger() *Logger {
    return &Logger{}
}

func (l *Logger) Record(decision Decision) {
    l.mu.Lock()
    defer l.mu.Unlock()
    fmt.Printf("Audit: %+v\n", decision)
}
