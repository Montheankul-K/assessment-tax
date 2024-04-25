package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type IConfig interface {
	App() IAppConfig
	DB() IDBConfig
	AdminAuth() IAdminAuth
}

type config struct {
	app       *app
	db        *db
	adminAuth *adminAuth
}

type IAppConfig interface {
	Port() string
}

type app struct {
	port string
}

type IDBConfig interface {
	Url() string
}

type db struct {
	url string
}

type IAdminAuth interface {
	Username() string
	Password() string
}

type adminAuth struct {
	username string
	password string
}

func (c *config) App() IAppConfig {
	return c.app
}

func (c *config) DB() IDBConfig {
	return c.db
}

func (c *config) AdminAuth() IAdminAuth {
	return c.adminAuth
}

func (a *app) Port() string {
	return a.port
}

func (d *db) Url() string {
	return d.url
}

func (a *adminAuth) Username() string {
	return a.username
}

func (a *adminAuth) Password() string {
	return a.password
}

func loadEnv(path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func LoadConfig(path string, useEnv bool) (IConfig, error) {
	if useEnv {
		loadEnv(path)
	}

	appPort := os.Getenv("PORT")
	if appPort == "" {
		return nil, fmt.Errorf("env variable PORT not set")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("env variable DATABASE_URL not set")
	}

	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		return nil, fmt.Errorf("env variable ADMIN_USERNAME not set")
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		return nil, fmt.Errorf("env variable ADMIN_PASSWORD not set")
	}

	return &config{
		app: &app{
			port: appPort,
		},
		db: &db{
			url: dbUrl,
		},
		adminAuth: &adminAuth{
			username: adminUsername,
			password: adminPassword,
		},
	}, nil
}
