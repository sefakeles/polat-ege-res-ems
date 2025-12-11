package health

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusDegraded  Status = "degraded"
)

type CheckResult struct {
	Name      string        `json:"name"`
	Status    Status        `json:"status"`
	Message   string        `json:"message,omitempty"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
}

type Checker interface {
	Name() string
	Check(ctx context.Context) error
}

type HealthService struct {
	checkers []Checker
	mutex    sync.RWMutex
	timeout  time.Duration
}

func NewHealthService() *HealthService {
	return &HealthService{
		checkers: make([]Checker, 0),
		timeout:  5 * time.Second,
	}
}

func (h *HealthService) RegisterChecker(checker Checker) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.checkers = append(h.checkers, checker)
}

func (h *HealthService) CheckAll(ctx context.Context) map[string]CheckResult {
	h.mutex.RLock()
	checkers := make([]Checker, len(h.checkers))
	copy(checkers, h.checkers)
	h.mutex.RUnlock()

	results := make(map[string]CheckResult)
	resultChan := make(chan CheckResult, len(checkers))

	// Run all checks concurrently
	for _, checker := range checkers {
		go func(c Checker) {
			result := h.runSingleCheck(ctx, c)
			resultChan <- result
		}(checker)
	}

	// Collect results
	for range checkers {
		result := <-resultChan
		results[result.Name] = result
	}

	return results
}

func (h *HealthService) runSingleCheck(ctx context.Context, checker Checker) CheckResult {
	start := time.Now()

	// Create context with timeout
	checkCtx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	result := CheckResult{
		Name:      checker.Name(),
		Timestamp: start,
	}

	// Run the check
	if err := checker.Check(checkCtx); err != nil {
		result.Status = StatusUnhealthy
		result.Message = err.Error()
	} else {
		result.Status = StatusHealthy
	}

	result.Duration = time.Since(start)
	return result
}

func (h *HealthService) GetOverallStatus(results map[string]CheckResult) Status {
	healthyCount := 0
	totalCount := len(results)

	for _, result := range results {
		if result.Status == StatusHealthy {
			healthyCount++
		}
	}

	if healthyCount == totalCount {
		return StatusHealthy
	} else if healthyCount > 0 {
		return StatusDegraded
	} else {
		return StatusUnhealthy
	}
}

// Database Health Checker
type DatabaseChecker struct {
	name string
	db   interface{ HealthCheck() error }
}

func NewDatabaseChecker(name string, db interface{ HealthCheck() error }) *DatabaseChecker {
	return &DatabaseChecker{name: name, db: db}
}

func (d *DatabaseChecker) Name() string {
	return d.name
}

func (d *DatabaseChecker) Check(ctx context.Context) error {
	return d.db.HealthCheck()
}

// Service Health Checker
type ServiceChecker struct {
	name    string
	service interface{ IsConnected() bool }
}

func NewServiceChecker(name string, service interface{ IsConnected() bool }) *ServiceChecker {
	return &ServiceChecker{name: name, service: service}
}

func (s *ServiceChecker) Name() string {
	return s.name
}

func (s *ServiceChecker) Check(ctx context.Context) error {
	if !s.service.IsConnected() {
		return fmt.Errorf("%s is not connected", s.name)
	}
	return nil
}
