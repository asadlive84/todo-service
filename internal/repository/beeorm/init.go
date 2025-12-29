package beeorm

import (
	"git.ice.global/packages/beeorm/v4"
	"git.ice.global/packages/hitrix/pkg/entity"
	e "todo-service/internal/repository/beeorm/entity"
)

func Init(registry *beeorm.Registry) {
	registry.RegisterEntity(
		&e.TodoEntity{},
	)
	registry.RegisterEntity(&entity.RequestLoggerEntity{})

}
