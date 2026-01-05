package beeorm

import (
	e "todo-service/internal/repository/beeorm/entity"

	"git.ice.global/packages/beeorm/v4"
	// "git.ice.global/packages/hitrix/service"
)

// func Init(registry *beeorm.Registry) {
// 	configService := service.DI().Config()
// 	REDIS_ADDR := configService.DefString("REDIS.REDIS_ADDR", "localhost:6379")

// 	registry.RegisterEntity(
// 		&e.TodoEntity{},
// 		&e.OutboxEntity{},
// 		&e.FileEnity{},
// 	)

// 	// registry.RegisterRedisStream("", "todo:events", []string{"todo-consumer-group"})

// 	registry.RegisterEntity(&entity.RequestLoggerEntity{})

// 	registry.RegisterRedis(REDIS_ADDR, "", 0, "todo_cache")
// 	registry.RegisterRedis(REDIS_ADDR, "", 0, "file_cache")

// 	registry.RegisterRedis(REDIS_ADDR, "", 0, "todo_search")
// 	registry.RegisterRedis(REDIS_ADDR, "", 0, "file_search")

// 	// registry.RegisterRedisStream("todo.events", "todo_cache", nil)
// 	registry.RegisterRedisStream("todo:events", "todo_cache", []string{"todo-consumer-group"})
// 	registry.RegisterRedisStream("files:events", "file_cache", nil)

// 	log.Println("BeeORM initialized successfully")
// }

func Init(registry *beeorm.Registry) {
	// configService := service.DI().Config()
	// REDIS_ADDR := configService.DefString("REDIS.REDIS_ADDR", "localhost:6379")
	REDIS_ADDR := "localhost:6379"

	registry.RegisterEntity(&e.TodoEntity{}, &e.OutboxEntity{}, &e.FileEnity{})
	// registry.RegisterEntity(&entity.RequestLoggerEntity{})

	registry.RegisterRedis(REDIS_ADDR, "", 0, "todo_cache")
	registry.RegisterRedis(REDIS_ADDR, "", 0, "file_cache")

	registry.RegisterRedisStream("todo:events", "todo_cache", []string{"todo-consumer-group"})
	registry.RegisterRedisStream("files:events", "file_cache", nil)
}
