package cache

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/estrategiahq/backend-libs/logger"
	"github.com/go-redis/redis"
)

var (
	connectionError = errors.New("cache: não foi possível realizar a conexão com o cache")

	setError = errors.New("cache: não foi possível escrever o objeto no cache")
	getError = errors.New("cache: não foi possível recuperar o objeto do cache")
)

// ICache interface for a new instance of cache
type ICache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, value interface{}) error
}

// Config model of configuration
type Config struct {
	Addr string
	DB   int
}

const defaultAddr = "localhost:6379"

// New function return a new instance
func New(config Config) ICache {
	addr := config.Addr
	db := config.DB

	if addr == "" {
		addr = defaultAddr
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       db,
	})

	return &cacheImpl{
		client: client,
	}
}

type cacheImpl struct {
	client *redis.Client
}

func (c cacheImpl) Set(ctx context.Context, key string, document interface{}) error {
	log := logger.FromContext(ctx)

	bt, err := json.Marshal(document)
	if err != nil {
		log.Error(setError, err.Error())
		return setError
	}

	err = c.client.Set(key, string(bt), 0).Err()
	if err != nil {
		log.Error(setError, err.Error())
		return setError
	}

	return nil
}

func (c cacheImpl) Get(ctx context.Context, key string, document interface{}) error {
	log := logger.FromContext(ctx)

	val, err := c.client.Get(key).Result()
	if err != nil {
		log.Error(setError, err.Error())
		return setError
	}

	err = json.Unmarshal([]byte(val), document)
	if err != nil {
		log.Error(setError, err.Error())
		return setError
	}

	return nil
}
