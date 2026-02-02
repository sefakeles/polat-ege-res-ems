package alarm

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"powerkonnekt/ems/internal/database"
)

// Module provides alarm management functionality to the Fx application
var Module = fx.Module("alarm",
	fx.Provide(ProvideManager),
)

// ProvideManager creates and provides an alarm manager instance
func ProvideManager(postgresDB *database.PostgresDB, logger *zap.Logger) *Manager {
	return NewManager(postgresDB, logger)
}
