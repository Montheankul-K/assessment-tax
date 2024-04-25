package servers

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/montheankul-k/assessment-tax/config"
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

func (s *server) setBasicAuth() {
	auth := s.config.AdminAuth()
	s.app.Use(middleware.BasicAuth(func(username string, password string, c echo.Context) (bool, error) {
		if username == auth.Username() && password == auth.Password() {
			return true, nil
		}

		return false, nil
	}))
}

func (s *server) setLogger() {
	s.app.Use(middleware.Logger())
}

func (s *server) setRecover() {
	s.app.Use(middleware.Recover())
}

func (s *server) setMiddleware() {
	s.setLogger()
	s.setRecover()
	s.setBasicAuth()
}

func (s *server) Start() {
	s.setMiddleware()

	port := s.config.App().Port()
	log.Println("Server started at port: " + port)

	go func() {
		if err := s.app.Start(":" + port); err != nil {
			log.Println("Starting server error: ", err)
		}
	}()

	s.gracefulShutdown()
}

func (s *server) gracefulShutdown() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.app.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error: ", err)
	}
}
