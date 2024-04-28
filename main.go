package main

import (
	"github.com/Montheankul-K/assessment-tax/config"
	"github.com/Montheankul-K/assessment-tax/modules/server"
	"github.com/Montheankul-K/assessment-tax/modules/tax"
	"github.com/Montheankul-K/assessment-tax/packages/database"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".env", false)
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	db := database.DBConnect(cfg.DB())
	err = db.AutoMigrate(&tax.TaxAllowance{}, &tax.TaxLevel{})
	if err != nil {
		log.Fatal("Error migrate database tables: ", err)
	}

	server.NewServer(cfg, db).Start()
}
