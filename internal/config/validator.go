package config

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// NewValidator creates a new validator with custom validations registered
func NewValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	// Register custom validation for aligned intervals
	if err := v.RegisterValidation("aligned_interval", validateAlignedInterval); err != nil {
		panic(fmt.Sprintf("failed to register custom validator: %v", err))
	}
	// Register custom validation for log paths
	if err := v.RegisterValidation("logpath", validateLogPath); err != nil {
		panic(fmt.Sprintf("failed to register custom validator: %v", err))
	}
	return v
}

// validateAlignedInterval validates that a duration aligns with time boundaries
func validateAlignedInterval(fl validator.FieldLevel) bool {
	interval, ok := fl.Field().Interface().(time.Duration)
	if !ok {
		return false
	}

	validIntervals := []time.Duration{
		// Sub-second intervals
		5 * time.Millisecond,
		10 * time.Millisecond,
		20 * time.Millisecond,
		25 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		200 * time.Millisecond,
		250 * time.Millisecond,
		500 * time.Millisecond,
		// Second intervals
		time.Second,
		2 * time.Second,
		5 * time.Second,
		10 * time.Second,
		15 * time.Second,
		20 * time.Second,
		30 * time.Second,
		// Minute intervals
		time.Minute,
		2 * time.Minute,
		3 * time.Minute,
		4 * time.Minute,
		5 * time.Minute,
		6 * time.Minute,
		10 * time.Minute,
		12 * time.Minute,
		15 * time.Minute,
		20 * time.Minute,
		30 * time.Minute,
		// Hour intervals
		time.Hour,
	}

	return slices.Contains(validIntervals, interval)
}

// validateLogPath validates that a log path is either stdout/stderr or a valid file path
func validateLogPath(fl validator.FieldLevel) bool {
	path := fl.Field().String()

	// Allow special output streams
	if path == "stdout" || path == "stderr" {
		return true
	}

	// Validate file paths - ensure not empty and trimmed
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}

	// Additional validation: ensure path doesn't contain only whitespace
	// and has at least some valid path characters
	if len(path) == 0 {
		return false
	}

	return true
}
