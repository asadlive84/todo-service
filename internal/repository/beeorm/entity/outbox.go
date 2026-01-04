package entity

import (
	"time"

	"git.ice.global/packages/beeorm/v4"
)

type OutboxEntity struct {
	beeorm.ORM  `orm:"table=outbox"`
	ID          uint64
	EventType   string
	Payload     string
	Status      string
	CreatedAt   time.Time `orm:"time"`
	PublishedAt time.Time `orm:"time"`
}
