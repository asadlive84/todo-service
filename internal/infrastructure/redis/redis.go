package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitRedis(address string) *redis.Client {
	fmt.Println("============adress redis===>", address)
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// ✅ Create stream (XADD creates stream if not exists)
	if err := createStreams(ctx, client); err != nil {
		log.Printf("Warning: Failed to initialize streams: %v", err)
	}

	log.Println("Redis connected and streams initialized")
	return client
}

func createStreams(ctx context.Context, client *redis.Client) error {
	streams := []string{
		"todos:events",
		"files:events",
	}

	for _, stream := range streams {
		// XADD with MAXLEN 0 creates stream without adding real data
		// Or just add a dummy message that you trim later
		err := client.XAdd(ctx, &redis.XAddArgs{
			Stream: stream,
			MaxLen: 1000, // Keep only last 1000 events
			Values: map[string]interface{}{
				"event":     "stream.initialized",
				"timestamp": "init",
			},
		}).Err()

		if err != nil {
			return err
		}

		log.Printf("Stream '%s' initialized", stream)
	}

	return nil
}
