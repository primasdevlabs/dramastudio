package metrics

import (
	"context"
	"sync"
)

// Collector records counters and durations under production metrics names.
type Collector interface {
	IncCounter(name string, value float64)
	RecordDuration(name string, seconds float64)
}

// MemoryCollector is a simple in-process collector suitable for dev/tests.
// Swap for an OTel meter-backed implementation when exporters are wired.
type MemoryCollector struct {
	mu       sync.Mutex
	counters map[string]float64
	timings  map[string][]float64
}

func NewMemory() *MemoryCollector {
	return &MemoryCollector{
		counters: make(map[string]float64),
		timings:  make(map[string][]float64),
	}
}

func (m *MemoryCollector) IncCounter(name string, value float64) {
	m.mu.Lock()
	m.counters[name] += value
	m.mu.Unlock()
}

func (m *MemoryCollector) RecordDuration(name string, seconds float64) {
	m.mu.Lock()
	m.timings[name] = append(m.timings[name], seconds)
	m.mu.Unlock()
}

func (m *MemoryCollector) Counter(name string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.counters[name]
}

// Noop drops all samples.
type Noop struct{}

func (Noop) IncCounter(string, float64)     {}
func (Noop) RecordDuration(string, float64) {}

// Ensure interface conformance.
var _ = context.Background
