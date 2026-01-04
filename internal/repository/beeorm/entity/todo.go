package entity

import (
	"time"

	"git.ice.global/packages/beeorm/v4"
)

type TodoEntity struct {
	// beeorm.ORM  `orm:"table=todo;redis=todo_cache;redisCache;redisSearch=todo_search"`
	beeorm.ORM  `orm:"table=todo;redis=todo_cache;redisCache;redisSearch=todo_search;dirty=todo:events"`
	ID          uint64    `orm:"pk;searchable;sortable"`
	Description string    `orm:"searchable"`
	DueDate     time.Time `orm:"time"`
	CreatedAt   time.Time `orm:"time"`
	FileID      string
}
