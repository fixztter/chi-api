package main

import (
	"log"

	"github.com/fixztter/chi-api/cmd/api"
	"github.com/fixztter/chi-api/config"
	"github.com/fixztter/chi-api/db"
)

func main() {
	db, err := db.NewPostgresStorage(config.Envs.ConnectionString)
	if err != nil {
		log.Fatalf("db: %s\n", err.Error())
	}

	s := api.NewAPIServer(config.Envs.Port, db)
	if err := s.Run(); err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}
