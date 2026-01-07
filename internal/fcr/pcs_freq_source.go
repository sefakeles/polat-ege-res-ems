package fcr

import (
	"fmt"
	"sync"
	"time"

	"powerkonnekt/ems/internal/pcs"
	"powerkonnekt/ems/pkg/logger"
)

// PCSFrequencySource implements FrequencySource using BESS PCS data
type PCSFrequencySource struct {
	pcsManager *pcs.Manager
	pcsNo      uint8
	log        logger.Logger

	// Frequency data
	mutex         sync.RWMutex
	lastFrequency float64
	lastUpdate    time.Time

	// Subscribers
	subscribers []func(float64)
	subMutex    sync.RWMutex
}

// NewPCSFrequencySource creates a new PCS-based frequency source
func NewPCSFrequencySource(pcsManager *pcs.Manager, pcsNo uint8) *PCSFrequencySource {
	sourceLogger := logger.With(
		logger.String("component", "pcs_frequency_source"),
		logger.Uint8("pcs_no", pcsNo),
	)

	return &PCSFrequencySource{
		pcsManager:  pcsManager,
		pcsNo:       pcsNo,
		log:         sourceLogger,
		subscribers: make([]func(float64), 0),
	}
}

// GetFrequency returns the latest grid frequency from PCS
func (p *PCSFrequencySource) GetFrequency() (float64, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// Check if data is fresh (within last 5 seconds)
	if time.Since(p.lastUpdate) > 5*time.Second {
		return 0, fmt.Errorf("frequency data is stale (last update: %v)", p.lastUpdate)
	}

	return p.lastFrequency, nil
}

// Subscribe adds a callback for frequency updates
func (p *PCSFrequencySource) Subscribe(callback func(float64)) error {
	p.subMutex.Lock()
	defer p.subMutex.Unlock()

	p.subscribers = append(p.subscribers, callback)
	return nil
}

// UpdateFromPCS reads frequency from PCS data and notifies subscribers
func (p *PCSFrequencySource) UpdateFromPCS() error {
	// Get PCS data from BESS service
	pcsService, err := p.pcsManager.GetService(int(p.pcsNo))
	if err != nil {
		return fmt.Errorf("failed to get PCS service %d: %w", p.pcsNo, err)
	}
	pcsData := pcsService.GetLatestPCSGridData()

	frequency := float64(pcsData.GridFrequency)

	// Update stored frequency
	p.mutex.Lock()
	p.lastFrequency = frequency
	p.lastUpdate = time.Now()
	p.mutex.Unlock()

	// Notify subscribers
	p.notifySubscribers(frequency)

	return nil
}

// notifySubscribers calls all registered callbacks with new frequency
func (p *PCSFrequencySource) notifySubscribers(frequency float64) {
	p.subMutex.RLock()
	defer p.subMutex.RUnlock()

	for _, callback := range p.subscribers {
		go callback(frequency)
	}
}
