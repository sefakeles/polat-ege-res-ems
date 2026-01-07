package fcr

import (
	"sync"
	"time"

	"powerkonnekt/ems/pkg/logger"
)

// TestFrequencySource implements FrequencySource with settable frequency
// for testing purposes
type TestFrequencySource struct {
	log logger.Logger

	// Frequency data
	mutex      sync.RWMutex
	frequency  float64
	lastUpdate time.Time

	// Subscribers
	subscribers []func(float64)
	subMutex    sync.RWMutex
}

// NewTestFrequencySource creates a new test frequency source
func NewTestFrequencySource() *TestFrequencySource {
	sourceLogger := logger.With(
		logger.String("component", "test_frequency_source"),
	)

	return &TestFrequencySource{
		log:         sourceLogger,
		frequency:   50.0, // Start at nominal frequency
		lastUpdate:  time.Now(),
		subscribers: make([]func(float64), 0),
	}
}

// GetFrequency returns the current test frequency
func (t *TestFrequencySource) GetFrequency() (float64, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.frequency, nil
}

// SetFrequency sets a new frequency value and notifies subscribers
func (t *TestFrequencySource) SetFrequency(frequency float64) {
	t.mutex.Lock()
	t.frequency = frequency
	t.lastUpdate = time.Now()
	t.mutex.Unlock()

	t.log.Info("Test frequency updated",
		logger.Float64("frequency", frequency))

	// Notify subscribers
	t.notifySubscribers(frequency)
}

// Subscribe adds a callback for frequency updates
func (t *TestFrequencySource) Subscribe(callback func(float64)) error {
	t.subMutex.Lock()
	defer t.subMutex.Unlock()

	t.subscribers = append(t.subscribers, callback)
	return nil
}

// notifySubscribers calls all registered callbacks with new frequency
func (t *TestFrequencySource) notifySubscribers(frequency float64) {
	t.subMutex.RLock()
	defer t.subMutex.RUnlock()

	for _, callback := range t.subscribers {
		go callback(frequency)
	}
}

// GetLastUpdate returns the time of last frequency update
func (t *TestFrequencySource) GetLastUpdate() time.Time {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.lastUpdate
}
