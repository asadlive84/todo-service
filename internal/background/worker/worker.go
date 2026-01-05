package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	beeOrmEnity "todo-service/internal/repository/beeorm/entity"

	"git.ice.global/packages/beeorm/v4"
)

type OutboxProcessor struct {
	engine *beeorm.Engine
}

func NewOutboxProcessor(engine *beeorm.Engine) *OutboxProcessor {
	return &OutboxProcessor{
		engine: engine,
	}
}

func (p *OutboxProcessor) Start(ctx context.Context) {
	log.Println(" Outbox processor started...")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println(" Outbox processor stopped")
			return

		case <-ticker.C:
			p.processEvents(ctx)
		}
	}
}

func (p *OutboxProcessor) processEvents(ctx context.Context) {
	// Context check
	if err := ctx.Err(); err != nil {
		log.Printf(" Context cancelled: %v", err)
		return
	}

	// Step 1: Find pending events
	var events []*beeOrmEnity.OutboxEntity
	where := beeorm.NewWhere("Status IN (?, ?)", "pending", "failed")
	pager := beeorm.NewPager(1, 50)

	p.engine.Search(where, pager, &events)

	if len(events) == 0 {
		log.Println("📭 No pending events")
		return
	}

	log.Printf(" Processing %d pending events...", len(events))

	// Process each event
	for _, event := range events {
		// Context check per event
		if err := ctx.Err(); err != nil {
			log.Printf("Context cancelled, stopping processing")
			return
		}

		if err := p.publishEvent(ctx, event); err != nil {
			log.Printf(" Failed to publish event %d: %v", event.ID, err)
			p.markAsFailed(event, err)
		} else {
			log.Printf("Published event %d (type: %s)", event.ID, event.EventType)
			p.markAsPublished(event)
		}
	}

	log.Printf("Completed processing %d events", len(events))
}

// Helper: Mark as failed
func (p *OutboxProcessor) markAsFailed(event *beeOrmEnity.OutboxEntity, err error) {
	event.Status = "failed"
	// event.RetryCount++
	// event.LastError = err.Error()

	if flushErr := p.engine.FlushWithCheck(event); flushErr != nil {
		log.Printf("Failed to update event %d status: %v", event.ID, flushErr)
	}
}

// Helper: Mark as published
func (p *OutboxProcessor) markAsPublished(event *beeOrmEnity.OutboxEntity) {
	now := time.Now()
	event.Status = "published"
	event.PublishedAt = now

	if err := p.engine.FlushWithCheck(event); err != nil {
		log.Printf("Failed to update event %d status: %v", event.ID, err)
	}
}

func (p *OutboxProcessor) publishEvent(ctx context.Context, event *beeOrmEnity.OutboxEntity) (err error) {
	// 	Context validation
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	// Input validation
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}
	if event.EventType == "" {
		return fmt.Errorf("event type is required")
	}
	if event.Payload == "" {
		return fmt.Errorf("event payload is required")
	}

	// Timeout context
	publishCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Recovery to catch BeeORM internal panics
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("beeorm panic: %v", r)
		}
	}()

	broker := p.engine.GetEventBroker()
	flusher := broker.NewFlusher()

	// 1. Queue the event
	flusher.Publish("todo:events", map[string]interface{}{
		"id":         fmt.Sprintf("%d", event.ID),
		"event_type": event.EventType,
		"payload":    event.Payload,
		"timestamp":  time.Now().Unix(),
		"created_at": event.CreatedAt.Unix(),
	})

	event.Status = "published"
	event.PublishedAt = time.Now()

	dbFlusher := p.engine.NewFlusher()
	dbFlusher.Track(event)

	select {
	case <-publishCtx.Done():
		return fmt.Errorf("publish timeout: %w", publishCtx.Err())
	default:
		dbFlusher.Flush()
		flusher.Flush()
	}

	return nil
}
