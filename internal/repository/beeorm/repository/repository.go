package repository

import (
	"git.ice.global/packages/beeorm/v4"
)

type OrmEngine struct {
	orm *beeorm.Engine
}

func NewOrmEngine(orm *beeorm.Engine) *OrmEngine {
	return &OrmEngine{orm: orm}
}
