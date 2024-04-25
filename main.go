package main

import (
	"github.com/KKGo-Software-engineering/assessment-tax/config"
	"github.com/KKGo-Software-engineering/assessment-tax/modules/servers"
	"github.com/KKGo-Software-engineering/assessment-tax/packages/database"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".env", true)
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	db := database.DBConnect(cfg.DB())
	servers.NewServer(cfg, db).Start()
}
