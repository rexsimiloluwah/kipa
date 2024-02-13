package queue

import (
	"fmt"
	"keeper/internal/config"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type RedisQueue struct {
	client    *asynq.Client
	inspector *asynq.Inspector
	opts      *asynq.RedisClientOpt
}

func NewRedisQueue(cfg *config.Config) *RedisQueue {
	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	redisConnectionOpts := asynq.RedisClientOpt{
		Addr: redisAddr,
	}
	client := asynq.NewClient(redisConnectionOpts)
	inspector := asynq.NewInspector(redisConnectionOpts)
	return &RedisQueue{
		client:    client,
		inspector: inspector,
		opts:      &redisConnectionOpts,
	}
}

func (q *RedisQueue) Inspector() *asynq.Inspector {
	return q.inspector
}

func (q *RedisQueue) Add(task *asynq.Task, opt asynq.Option) error {
	info, err := q.client.Enqueue(task, opt)
	if err != nil {
		logrus.Info("enqueued task: id=%s, queue=%s", info.ID, info.Queue)
		return err
	}
	return nil
}
