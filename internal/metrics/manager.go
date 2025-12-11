package metrics

import (
	"context"
	"runtime"
	"sync"
	"time"

	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/logger"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/net"
)

// Manager handles metrics collection and storage
type Manager struct {
	db     *database.InfluxDB
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mutex  sync.RWMutex
	log    logger.Logger

	// Runtime metrics
	startTime time.Time

	// Network metrics
	lastNetRx uint64
	lastNetTx uint64
}

// NewManager creates a new metrics manager
func NewManager(db *database.InfluxDB) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	// Create component-specific logger
	metricsLogger := logger.With(
		logger.String("component", "metrics_manager"),
	)

	return &Manager{
		db:        db,
		ctx:       ctx,
		cancel:    cancel,
		startTime: time.Now(),
		log:       metricsLogger,
	}
}

// Start starts the metrics collection
func (m *Manager) Start() error {
	// Initialize network counters
	m.initNetworkCounters()

	m.wg.Go(m.collectLoop)

	m.log.Info("Metrics manager started",
		logger.Time("start_time", m.startTime))
	return nil
}

// Stop stops the metrics collection
func (m *Manager) Stop() {
	m.cancel()
	m.wg.Wait()
	m.log.Info("Metrics manager stopped")
}

func (m *Manager) collectLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.collectSystemMetrics()
			m.collectRuntimeMetrics()
		}
	}
}

func (m *Manager) initNetworkCounters() {
	netStats, err := net.IOCounters(false)
	if err != nil || len(netStats) == 0 {
		m.log.Error("Failed to initialize network counters", logger.Err(err))
		return
	}

	m.mutex.Lock()
	m.lastNetRx = netStats[0].BytesRecv
	m.lastNetTx = netStats[0].BytesSent
	m.mutex.Unlock()
}

func (m *Manager) collectSystemMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Get CPU usage
	cpuPercent, err := cpu.Percent(time.Second, false)
	var cpuUsage float32
	if err != nil || len(cpuPercent) == 0 {
		m.log.Error("Failed to get CPU usage", logger.Err(err))
		cpuUsage = 0.0
	} else {
		cpuUsage = float32(cpuPercent[0])
	}

	// Get disk usage for root partition
	diskStat, err := disk.Usage("/")
	var diskUsage float32
	if err != nil {
		m.log.Error("Failed to get disk usage", logger.Err(err))
		diskUsage = 0.0
	} else {
		diskUsage = float32(diskStat.UsedPercent)
	}

	// Get network statistics
	netRx, netTx := m.getNetworkStats()

	metrics := database.SystemMetrics{
		Timestamp:   time.Now(),
		CPUUsage:    cpuUsage,
		MemoryUsage: float32(memStats.Alloc) / 1024 / 1024, // MB
		DiskUsage:   diskUsage,
		NetworkRx:   netRx,
		NetworkTx:   netTx,
	}

	// Save to InfluxDB
	if err := m.db.WriteSystemMetrics(metrics); err != nil {
		m.log.Error("Failed to save system metrics to InfluxDB", logger.Err(err))
	}
}

func (m *Manager) getNetworkStats() (uint64, uint64) {
	netStats, err := net.IOCounters(false)
	if err != nil || len(netStats) == 0 {
		m.log.Error("Failed to get network statistics", logger.Err(err))
		return 0, 0
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	currentRx := netStats[0].BytesRecv
	currentTx := netStats[0].BytesSent

	// Calculate delta since last measurement
	deltaRx := currentRx - m.lastNetRx
	deltaTx := currentTx - m.lastNetTx

	// Update last values
	m.lastNetRx = currentRx
	m.lastNetTx = currentTx

	return deltaRx, deltaTx
}

// collectRuntimeMetrics collects and stores runtime metrics
func (m *Manager) collectRuntimeMetrics() {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	runtimeMetrics := database.RuntimeMetrics{
		Timestamp:      time.Now(),
		UptimeSeconds:  time.Since(m.startTime).Seconds(),
		Goroutines:     runtime.NumGoroutine(),
		HeapAllocMB:    float64(memStats.HeapAlloc) / 1024 / 1024,
		HeapSysMB:      float64(memStats.HeapSys) / 1024 / 1024,
		HeapIdleMB:     float64(memStats.HeapIdle) / 1024 / 1024,
		HeapInUseMB:    float64(memStats.HeapInuse) / 1024 / 1024,
		HeapReleasedMB: float64(memStats.HeapReleased) / 1024 / 1024,
		StackInUseMB:   float64(memStats.StackInuse) / 1024 / 1024,
		StackSysMB:     float64(memStats.StackSys) / 1024 / 1024,
		GCRuns:         uint32(memStats.NumGC),
		GCPauseTotalNs: memStats.PauseTotalNs,
		GCCPUFraction:  memStats.GCCPUFraction,
		NextGCMB:       float64(memStats.NextGC) / 1024 / 1024,
		LastGCTime:     time.Unix(0, int64(memStats.LastGC)).Unix(),
		MallocsTotal:   memStats.Mallocs,
		FreesTotal:     memStats.Frees,
		TotalAllocMB:   float64(memStats.TotalAlloc) / 1024 / 1024,
		LookupsTotal:   memStats.Lookups,
	}

	// Save to InfluxDB
	if err := m.db.WriteRuntimeMetrics(runtimeMetrics); err != nil {
		m.log.Error("Failed to save runtime metrics to InfluxDB", logger.Err(err))
	}
}
