package alarm

import (
	"go.uber.org/fx"

	"powerkonnekt/ems/internal/database"
)

// Module provides alarm management functionality to the Fx application
var Module = fx.Module("alarm",
	fx.Provide(ProvideManager),
)

// ProvideManager creates and provides an alarm manager instance
func ProvideManager(db *database.PostgresDB) *Manager {
	return NewManager(db)
}
