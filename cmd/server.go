package main

import (
	"keeper/internal/config"
	"keeper/internal/server"
	"keeper/pkg/mongo"
)

func main() {
	cfg := config.New()
	db := mongo.NewConnection(cfg)
	defer db.Disconnect()
	server := server.NewServer(cfg, db.Client)
	server.Start()
}
