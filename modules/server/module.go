package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/montheankul-k/assessment-tax/modules/admin/adminHandlers"
	"github.com/montheankul-k/assessment-tax/modules/middleware/middlewareHandlers"
	"github.com/montheankul-k/assessment-tax/modules/monitor/monitorHandlers"
	"github.com/montheankul-k/assessment-tax/modules/tax/taxHandlers"
	"github.com/montheankul-k/assessment-tax/modules/tax/taxRepositories"
	"github.com/montheankul-k/assessment-tax/modules/tax/taxUsecases"
)

type IModule interface {
	TaxModule()
	HealthCheckModule()
	AdminModule()
}

type moduleFactory struct {
	router     *echo.Echo
	server     *server
	middleware middlewareHandlers.IMiddlewareHandler
}

func NewModule(router *echo.Echo, server *server, middleware middlewareHandlers.IMiddlewareHandler) IModule {
	return &moduleFactory{
		router:     router,
		server:     server,
		middleware: middleware,
	}
}

func NewMiddleware(s *server) middlewareHandlers.IMiddlewareHandler {
	repository := taxRepositories.TaxRepository(s.db)
	usecase := taxUsecases.TaxUsecase(repository)
	handler := taxHandlers.TaxHandler(s.config, usecase)

	return middlewareHandlers.MiddlewareHandler(s.config, handler)
}

func (m *moduleFactory) basicAuthMiddleware(username, password string) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(user, pass string, c echo.Context) (bool, error) {
		if user == username && pass == password {
			return true, nil
		}

		return false, nil
	})
}

func (m *moduleFactory) HealthCheckModule() {
	handler := monitorHandlers.MonitorHandler(m.server.config)

	router := m.router.Group("")
	router.GET("/health", handler.HealthCheck)
}

func (m *moduleFactory) TaxModule() {
	repository := taxRepositories.TaxRepository(m.server.db)
	usecase := taxUsecases.TaxUsecase(repository)
	handler := taxHandlers.TaxHandler(m.server.config, usecase)

	router := m.router.Group("/tax")
	router.POST("/calculations", m.middleware.ValidateCalculateTaxRequest(handler.CalculateTax))
	router.POST("/calculations/upload-csv", handler.CalculateTax)
}

func (m *moduleFactory) AdminModule() {
	auth := m.server.config.AdminAuth()

	repository := taxRepositories.TaxRepository(m.server.db)
	usecase := taxUsecases.TaxUsecase(repository)
	handler := adminHandlers.AdminHandler(m.server.config, usecase)

	router := m.router.Group("/admin", m.basicAuthMiddleware(auth.Username(), auth.Password()))
	router.POST("/deductions/personal", m.middleware.ValidateSetDeductionRequest(handler.SetPersonalDeduction))
	router.POST("/deductions/k-receipt", m.middleware.ValidateSetDeductionRequest(handler.SetKReceiptDeduction))
}
