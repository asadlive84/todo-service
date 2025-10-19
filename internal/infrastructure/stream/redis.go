// todo-service/internal/infrastructure/stream/redis.go

package stream

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"todo-service/internal/domain/entity"
	_interface "todo-service/internal/domain/interface"
)

type redisStreamRepository struct {
	client *redis.Client
	stream string
}

func NewRedisStreamRepository(addr, stream string) _interface.RedisStreamRepository {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &redisStreamRepository{client: rdb, stream: stream}
}

func (r *redisStreamRepository) PublishTodo(ctx context.Context, todo *entity.TodoItem) error {
	data, err := json.Marshal(todo)
	if err != nil {
		return err
	}
	return r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: r.stream,
		Values: map[string]interface{}{"data": string(data)},
	}).Err()
}
