package entity

import (
	"time"

	"git.ice.global/packages/beeorm/v4"
)

//`orm:"table=todos;mysql=todo_pool;redis=cache;redisCache;redisSearch=search;dirty=beeorm.fid"`

type TodoEntity struct {
	beeorm.ORM  `orm:"table=todos;redis=todo_cache;redisCache;redisSearch=todo_search;dirty=todo.events"`
	ID          uint64    `orm:"pk;searchable;sortable"`
	Description string    `orm:"searchable;sortable"`
	DueDate     time.Time `orm:"time"`
	CreatedAt   time.Time `orm:"time"`
	FileID      string
}
