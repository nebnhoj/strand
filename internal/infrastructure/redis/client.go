package redis

import (
	"context"
	"log"
	"os"

	goredis "github.com/redis/go-redis/v9"
)

func NewClient() *goredis.Client {
	opt, err := goredis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("redis: parse URL: %v", err)
	}
	client := goredis.NewClient(opt)
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis: ping: %v", err)
	}
	return client
}
