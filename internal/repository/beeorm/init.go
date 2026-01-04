package beeorm

import (
	"log"
	e "todo-service/internal/repository/beeorm/entity"

	"git.ice.global/packages/beeorm/v4"
	"git.ice.global/packages/hitrix/pkg/entity"
	"git.ice.global/packages/hitrix/service"
)

// func Init(registry *beeorm.Registry) {
// 	configService := service.DI().Config()
// 	REDIS_ADDR := configService.DefString("REDIS.REDIS_ADDR", "localhost:6379")
// 	// REDIS_STREAM := configService.DefString("REDIS.REDIS_STREAM", "todos:events")

// 	registry.RegisterEntity(
// 		&e.TodoEntity{},
// 		&e.FileEnity{},
// 	)
// 	registry.RegisterEntity(&entity.RequestLoggerEntity{})

// 	registry.RegisterRedis(REDIS_ADDR, "", 0, "cache")

// 	registry.RegisterRedisStream("beeorm.fid", "cache", nil)

// 	// registry.RegisterRedis(REDIS_ADDR, "", 0, "search")

// 	log.Println("BeeORM initialized successfully")

// }

func Init(registry *beeorm.Registry) {
	configService := service.DI().Config()
	REDIS_ADDR := configService.DefString("REDIS.REDIS_ADDR", "localhost:6379")

	registry.RegisterEntity(
		&e.TodoEntity{},
		&e.FileEnity{},
	)

	registry.RegisterEntity(&entity.RequestLoggerEntity{})

	registry.RegisterRedis(REDIS_ADDR, "", 0, "todo_cache")
	registry.RegisterRedis(REDIS_ADDR, "", 0, "file_cache")

	registry.RegisterRedis(REDIS_ADDR, "", 0, "todo_search")
	registry.RegisterRedis(REDIS_ADDR, "", 0, "file_search")

	registry.RegisterRedisStream("todo.events", "todo_cache", nil)
	registry.RegisterRedisStream("files.events", "file_cache", nil)

	log.Println("BeeORM initialized successfully")
}
