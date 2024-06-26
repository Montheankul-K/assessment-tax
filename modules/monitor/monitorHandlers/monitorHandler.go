package monitorHandlers

import (
	"github.com/Montheankul-K/assessment-tax/config"
	"github.com/Montheankul-K/assessment-tax/modules/monitor"
	"github.com/labstack/echo/v4"
	"net/http"
)

type IMonitorHandler interface {
	HealthCheck(c echo.Context) error
}

type monitorHandler struct {
	config config.IConfig
}

func MonitorHandler(config config.IConfig) IMonitorHandler {
	return &monitorHandler{
		config: config,
	}
}

func (h *monitorHandler) HealthCheck(c echo.Context) error {
	res := &monitor.Monitor{
		Name:    h.config.App().Name(),
		Version: h.config.App().Version(),
	}

	return c.JSON(http.StatusOK, res)
}
