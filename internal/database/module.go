package database

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"powerkonnekt/ems/internal/config"
)

// Module provides database connections to the Fx application
var Module = fx.Module("database",
	fx.Provide(
		ProvideInfluxDB,
		ProvidePostgreSQL,
	),
	fx.Invoke(
		RegisterInfluxDBLifecycle,
		RegisterPostgreSQLLifecycle,
	),
)

// ProvideInfluxDB creates and provides an InfluxDB connection
func ProvideInfluxDB(cfg *config.Config, logger *zap.Logger) (*InfluxDB, error) {
	return NewInfluxDB(cfg.InfluxDB, logger)
}

// ProvidePostgreSQL creates and provides a PostgreSQL connection
func ProvidePostgreSQL(cfg *config.Config, logger *zap.Logger) (*PostgreSQL, error) {
	return NewPostgreSQL(cfg.PostgreSQL, logger)
}

// RegisterInfluxDBLifecycle registers lifecycle hooks for InfluxDB
func RegisterInfluxDBLifecycle(lc fx.Lifecycle, influxDB *InfluxDB) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return influxDB.Close()
		},
	})
}

// RegisterPostgreSQLLifecycle registers lifecycle hooks for PostgreSQL
func RegisterPostgreSQLLifecycle(lc fx.Lifecycle, postgreSQL *PostgreSQL) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return postgreSQL.Close()
		},
	})
}
