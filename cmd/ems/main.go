package main

import (
	"os"
	"os/signal"
	"syscall"

	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/ems"
	"powerkonnekt/ems/pkg/logger"
)

func main() {
	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Load configuration
	cfg, err := config.Load("configs/config.json")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Convert config.LoggerConfig to logger.Config
	loggerConfig := logger.Config{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
	}

	// Initialize logger
	if err := logger.InitializeWithConfig(loggerConfig); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Starting EMS application",
		logger.String("version", "dev"),
		logger.String("log_level", cfg.Logger.Level))

	// Initialize EMS
	emsInstance, err := ems.New(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize EMS", logger.Err(err))
	}

	// Start EMS
	if err := emsInstance.Start(); err != nil {
		logger.Fatal("Failed to start EMS", logger.Err(err))
	}

	// Wait for shutdown signal
	<-sigChan
	logger.Info("Shutdown signal received")

	// Stop EMS
	emsInstance.Stop()
	logger.Info("EMS application stopped")
}
