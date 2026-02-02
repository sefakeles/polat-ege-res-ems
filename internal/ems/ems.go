package ems

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/control"
)

// EMS represents the main EMS application
type EMS struct {
	config       config.EMSConfig
	controlLogic *control.Logic
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	log          *zap.Logger
}

// New creates a new EMS instance
func New(cfg config.EMSConfig, controlLogic *control.Logic, logger *zap.Logger) *EMS {
	ctx, cancel := context.WithCancel(context.Background())

	emsLogger := logger.With(
		zap.String("component", "ems"),
	)

	return &EMS{
		config:       cfg,
		controlLogic: controlLogic,
		ctx:          ctx,
		cancel:       cancel,
		log:          emsLogger,
	}
}

// Start starts the EMS
func (e *EMS) Start() error {
	e.wg.Go(e.reactiveControlLoop)
	e.log.Info("EMS application started")
	return nil
}

// Stop stops the EMS
func (e *EMS) Stop(ctx context.Context) {
	e.cancel()
	e.wg.Wait()
	e.log.Info("EMS application stopped")
}

// reactiveControlLoop runs reactive control logic triggered by data updates
func (e *EMS) reactiveControlLoop() {
	bessUpdateChan := e.controlLogic.GetBESSUpdateChannel()

	// Also run periodic control as a safety fallback
	fallbackTicker := time.NewTicker(100 * time.Millisecond)
	defer fallbackTicker.Stop()

	for {
		select {
		case <-e.ctx.Done():
			return
		case <-bessUpdateChan:
			// BESS data updated, execute control immediately
			// controlLogic.ExecuteControl()
		case <-fallbackTicker.C:
			// Safety fallback - ensure control runs at least once per 100 milliseconds
			e.controlLogic.ExecuteControl()
		}
	}
}
