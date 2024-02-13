package main

import (
	"keeper/internal/config"
	"keeper/internal/queue"
	"keeper/internal/queue/tasks"
	"keeper/internal/server"
	"keeper/pkg/mongo"
)

func main() {
	cfg := config.New()
	db := mongo.NewConnection(cfg)
	defer db.Disconnect()

	if cfg.WithWorkers {
		// register the asynq worker
		consumer := queue.NewConsumer(cfg)
		consumer.RegisterHandler(tasks.TypeUserVerificationMail, tasks.SendUserVerificationMail)
		consumer.RegisterHandler(tasks.TypeUserResetPasswordMail, tasks.SendResetPasswordMail)

		go consumer.Start()
	}

	server := server.NewServer(cfg, db.Client)
	server.Start()
}
