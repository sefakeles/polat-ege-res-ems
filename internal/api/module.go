package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/bms"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/control"
	"powerkonnekt/ems/internal/health"
	"powerkonnekt/ems/internal/pcs"
	"powerkonnekt/ems/internal/plc"
	"powerkonnekt/ems/internal/windfarm"
	"powerkonnekt/ems/pkg/logger"
)

// Module provides API server functionality to the Fx application
var Module = fx.Module("api",
	fx.Provide(
		ProvideHandlers,
		ProvideRouter,
		ProvideHTTPServer,
	),
	fx.Invoke(RegisterLifecycle),
)

// ProvideHandlers creates the API handlers
func ProvideHandlers(
	config *config.Config,
	bmsManager *bms.Manager,
	pcsManager *pcs.Manager,
	plcManager *plc.Manager,
	windFarmManager *windfarm.Manager,
	alarmManager *alarm.Manager,
	controlLogic *control.Logic,
	healthService *health.HealthService,
) *Handlers {
	return NewHandlers(
		config,
		bmsManager,
		pcsManager,
		plcManager,
		windFarmManager,
		alarmManager,
		controlLogic,
		healthService,
	)
}

// ProvideRouter creates and configures the Gin router
func ProvideRouter(handlers *Handlers) *gin.Engine {
	return SetupRoutes(handlers)
}

// ProvideHTTPServer creates the HTTP server
func ProvideHTTPServer(cfg *config.Config, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.EMS.HTTPPort),
		Handler: router,
	}
}

// RegisterLifecycle registers lifecycle hooks for the HTTP server
func RegisterLifecycle(lc fx.Lifecycle, server *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting HTTP server", logger.String("addr", server.Addr))
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("HTTP server error", logger.Err(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server")
			return server.Shutdown(ctx)
		},
	})
}
