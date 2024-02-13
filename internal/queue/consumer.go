package queue

import (
	"context"
	"fmt"
	"keeper/internal/config"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	mux *asynq.ServeMux
	srv *asynq.Server
}

func NewConsumer(cfg *config.Config) *Consumer {
	// create the redis connection
	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	redisConnection := asynq.RedisClientOpt{
		Addr: redisAddr,
	}

	// create and configure the asynq worker server
	server := asynq.NewServer(
		redisConnection,
		asynq.Config{
			Concurrency: 4, // number of concurrent workers to use
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			}, // specifying multiple queues with different priority
		},
	)

	mux := asynq.NewServeMux()

	return &Consumer{
		srv: server,
		mux: mux,
	}
}

func (c *Consumer) Start() {
	logrus.Info("starting asynq worker...")
	if err := c.srv.Run(c.mux); err != nil {
		logrus.WithError(err).Fatal("error starting asynq worker")
	}
}

func (c *Consumer) Stop() {
	logrus.Info("stopping asynq worker...")
	c.srv.Stop()
	c.srv.Shutdown()
}

func (c *Consumer) RegisterHandler(taskName string, handler func(context.Context, *asynq.Task) error) {
	c.mux.HandleFunc(taskName, handler)
}
