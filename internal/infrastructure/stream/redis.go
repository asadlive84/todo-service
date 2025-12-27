// todo-service/internal/infrastructure/stream/redis.go

package stream

import (
	"context"
	"encoding/json"
	"todo-service/internal/domain/entity"
	_interface "todo-service/internal/interface"

	"github.com/redis/go-redis/v9"
)

type redisStreamRepository struct {
	client *redis.Client
	stream string
}

func NewRedisStreamRepository(addr, stream string) _interface.RedisStreamRepository {
	
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &redisStreamRepository{client: rdb, stream: stream}
}

func (r *redisStreamRepository) PublishTodo(ctx context.Context, todo *entity.TodoItemEntity) error {
	data, err := json.Marshal(todo)
	if err != nil {
		return err
	}
	return r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: r.stream,
		Values: map[string]interface{}{"data": string(data)},
	}).Err()
}
