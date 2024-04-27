package main

import (
	"github.com/montheankul-k/assessment-tax/config"
	"github.com/montheankul-k/assessment-tax/modules/server"
	"github.com/montheankul-k/assessment-tax/modules/tax"
	"github.com/montheankul-k/assessment-tax/packages/database"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".env", true)
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
