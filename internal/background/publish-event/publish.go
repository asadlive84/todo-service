package publishevent

import (
	"context"
	"fmt"
	"log"
	"time"
	beeOrmEntity "todo-service/internal/repository/beeorm/entity"

	"git.ice.global/packages/beeorm/v4"
	"git.ice.global/packages/hitrix/service"
	"git.ice.global/packages/hitrix/service/component/app"
)

type PublishEvent struct{}

func (s *PublishEvent) Code() string {
	return "publish-event"
}

func (s *PublishEvent) Unique() bool {
	return false
}

func (s *PublishEvent) Description() string {
	return "Outbox pattern event processor"
}

func (s *PublishEvent) Infinity() bool {
	return true
}

// func (s *PublishEvent) Interval() time.Duration {
// 	return 2 * time.Second
// }

func (s *PublishEvent) Run(ctx context.Context, exit app.IExit) {

	time.Sleep(2 * time.Second)
	engine := service.DI().OrmEngine()
	processor := NewPublishEventProcess(engine)

	processor.ProcessOnce(ctx)
}

type PublishEventProcess struct {
	engine *beeorm.Engine
}

func NewPublishEventProcess(engine *beeorm.Engine) *PublishEventProcess {
	return &PublishEventProcess{engine: engine}
}

func (p *PublishEventProcess) ProcessOnce(ctx context.Context) {
	var events []*beeOrmEntity.OutboxEntity
	where := beeorm.NewWhere("Status IN (?, ?)", "pending", "failed")
	pager := beeorm.NewPager(1, 50)

	p.engine.Search(where, pager, &events)

	if len(events) == 0 {
		return
	}

	log.Printf("📦 Processing %d pending events...", len(events))

	var success, failed int

	for _, event := range events {
		if err := p.publishEvent(event); err != nil {
			log.Printf("Event %d failed: %v", event.ID, err)
			p.markAsFailed(event, err)
			failed++
		} else {
			log.Printf("Event %d published", event.ID)
			p.markAsPublished(event)
			success++
		}
	}

	log.Printf("Batch: %d success, %d failed", success, failed)
}

func (p *PublishEventProcess) publishEvent(event *beeOrmEntity.OutboxEntity) error {
	if event == nil || event.EventType == "" {
		return fmt.Errorf("invalid event data")
	}

	broker := p.engine.GetEventBroker()
	flusher := broker.NewFlusher()

	flusher.Publish("todo:events", map[string]interface{}{
		"id":         fmt.Sprintf("%d", event.ID),
		"event_type": event.EventType,
		"payload":    event.Payload,
		"timestamp":  time.Now().Unix(),
	})

	flusher.Flush()
	return nil
}

func (p *PublishEventProcess) markAsPublished(event *beeOrmEntity.OutboxEntity) {
	event.Status = "published"
	event.PublishedAt = time.Now()
	p.engine.FlushWithCheck(event)
}

func (p *PublishEventProcess) markAsFailed(event *beeOrmEntity.OutboxEntity, err error) {
	event.Status = "failed"
	p.engine.FlushWithCheck(event)
}
