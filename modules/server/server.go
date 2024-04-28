package server

import (
	"context"
	"github.com/Montheankul-K/assessment-tax/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type IServer interface {
	GetServer() *server
	Start()
}

type server struct {
	app    *echo.Echo
	config config.IConfig
	db     *gorm.DB
}

func NewServer(config config.IConfig, db *gorm.DB) IServer {
	return &server{
		app:    echo.New(),
		config: config,
		db:     db,
	}
}
func (s *server) GetServer() *server {
	return s
}

func (s *server) setLogger() {
	s.app.Use(middleware.Logger())
}

func (s *server) setRecover() {
	s.app.Use(middleware.Recover())
}

func (s *server) InitMiddleware() {
	s.setLogger()
	s.setRecover()
}

func (s *server) Start() {
	s.InitMiddleware()

	modules := NewModule(s.app, s, NewMiddleware(s))
	modules.HealthCheckModule()
	modules.TaxModule()
	modules.AdminModule()

	port := s.config.App().Port()
	log.Println("server started at port: " + port)

	go func() {
		if err := s.app.Start(":" + port); err != nil {
			log.Println("starting server error: ", err)
		}
	}()

	s.gracefulShutdown()
}

func (s *server) gracefulShutdown() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	log.Println("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.app.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown error: ", err)
	}
}
