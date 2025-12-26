package entity

import (
	"git.ice.global/packages/beeorm/v4"
	"git.ice.global/packages/hitrix/pkg/entity"
)

func Init(registry *beeorm.Registry) {
	registry.RegisterEntity(
		&TodoItemEntity{},
	)
	registry.RegisterEntity(&entity.RequestLoggerEntity{})

}
