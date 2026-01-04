package entity

import "time"

type OutboxEntity struct {
	ID          uint64
	EventType   string
	Payload     string
	Status      string
	CreatedAt   time.Time
	PublishedAt *time.Time
}
