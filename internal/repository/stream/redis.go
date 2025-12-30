
package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"todo-service/internal/domain/entity"

	"github.com/redis/go-redis/v9"
)

type RedisStreamRepository struct {
	client *redis.Client
	stream string
}

func NewRedisStreamRepository(addr, stream string) *RedisStreamRepository {

	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &RedisStreamRepository{client: rdb, stream: stream}
}

func (r *RedisStreamRepository) PublishTodo(ctx context.Context, todo *entity.TodoItem) error {

	fmt.Println("========PublishTodo=============")

	data, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	// return r.client.XAdd(ctx, &redis.XAddArgs{
	// 	Stream: "todos:events",
	// 	Values: map[string]interface{}{
	// 		"event":     "todos.created",
	// 		"data":      string(data),
	// 		"timestamp": todo.CreatedAt.Unix(),
	// 	},
	// }).Err()// this is not working

	return r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: r.stream,
		Values: map[string]interface{}{"data": string(data)},
	}).Err()
}
