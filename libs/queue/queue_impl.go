package queue

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/igorhalfeld/madalena-backend/lib/logger"
	"github.com/igorhalfeld/madalena-backend/lib/model"

	"github.com/go-redis/redis/v7"
)

var (
	errConnection = errors.New("queue: não foi possível realizar a conexão com a queue;")

	pubError = errors.New("queue: não foi possível publicar o evento;")
	subError = errors.New("queue: não foi possível subscrever no evento;")
)

// IQueue interface for a new instance of queue
type IQueue interface {
	Pub(ctx context.Context, channel string, value interface{}) error
	Sub(ctx context.Context, channel string) chan *model.Message
}

// Config model of configuration
type Config struct {
	DB       int
	URL      string
	Password string
}

const defaultURL = "localhost:6379"

// New function return a new instance
func New(config Config) IQueue {
	url := config.URL
	db := config.DB
	password := config.Password

	if url == "" {
		url = defaultURL
	}

	client := redis.NewClient(&redis.Options{
		DB:       db,
		Addr:     url,
		Password: password,
	})

	return &queueImpl{
		client: client,
	}
}

type queueImpl struct {
	client *redis.Client
}

func (c queueImpl) Pub(ctx context.Context, channel string, document interface{}) error {
	log := logger.FromContext(ctx)

	bt, err := json.Marshal(document)
	if err != nil {
		log.Error(pubError, err.Error())
		return pubError
	}

	err = c.client.Publish(channel, string(bt)).Err()
	if err != nil {
		log.Error(pubError, err.Error())
		return pubError
	}

	return nil
}

func (c queueImpl) Sub(ctx context.Context, channel string) chan *model.Message {
	sub := c.client.Subscribe(channel)
	ch := make(chan *model.Message)
	go func() {
		for {
			msg := <-sub.Channel()
			ch <- &model.Message{
				Channel: msg.Channel,
				Payload: msg.Payload,
			}
		}
	}()
	return ch
}
