package client

import (
	"sync"

	"github.com/go-redis/redis/v8"
)

var once sync.Once
var client *redis.Client

func Client() *redis.Client {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	})
	return client
}
