package redis

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"keeper/internal/config"

// 	"github.com/go-redis/redis/v8"
// 	r "github.com/go-redis/redis/v8"
// )

// type RedisClient struct {
// 	Client *r.Client
// }

// func NewRedisClient(cfg *config.Config) *RedisClient {
// 	var redisConnUri string
// 	if cfg.Env == "development" {
// 		redisConnUri = fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
// 	} else {
// 		redisConnUri = cfg.RedisProdUri
// 	}

// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     redisConnUri,
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})
// 	return &RedisClient{
// 		Client: rdb,
// 	}
// }

// func (r RedisClient) SetValueCache(data interface{}, key string) error {
// 	// marshal the data into bytes
// 	b, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}
// 	// set the value in the cache for that key
// 	err = r.Client.Set(context.Background(), key, b, 0).Err()
// 	return err
// }

// func (r RedisClient) GetKeyCache(key string) ([]byte, error) {
// 	value, err := r.Client.Get(context.Background(), key).Bytes()
// 	if err != nil {
// 		return nil, errors.New("key does not exist in cache")
// 	}
// 	return value, nil
// }
