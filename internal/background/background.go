package background

import (
	"context"
	"encoding/json"

	// "encoding/json"
	"fmt"
	"time"

	"git.ice.global/packages/beeorm/v4"
	"git.ice.global/packages/hitrix/service"
	"git.ice.global/packages/hitrix/service/component/app"
)

type EventConsumer struct{}

func (script *EventConsumer) Code() string {
	return "event-consumer"
}

func (script *EventConsumer) Unique() bool {
	return true
}

func (script *EventConsumer) Description() string {
	return "BeeORM event consumer"
}

func (script *EventConsumer) Infinity() bool {
	return true
}

func (script *EventConsumer) Run(ctx context.Context, exit app.IExit) {
	time.Sleep(2 * time.Second)
	fmt.Println("🚀 Starting event consumer...")

	engine := service.DI().OrmEngine()
	consumer := engine.GetEventBroker().Consumer("todo-consumer-group")

	fmt.Println("✅ Consumer ready, listening for events...")

	consumer.Consume(ctx, 1, func(events []beeorm.Event) {
		if len(events) == 0 {
			return
		}

		fmt.Printf("\n📦 Received EventConsumer %d events\n", len(events))
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		for _, event := range events {
			fmt.Printf("📦 Event ID: %s\n", event.ID())

			var fullPayload map[string]interface{}
			event.Unserialize(&fullPayload)

			payloadString, ok := fullPayload["payload"].(string)
			if !ok {
				fmt.Println("❌ Payload is not a string")
				continue
			}

			var data map[string]interface{}
			err := json.Unmarshal([]byte(payloadString), &data)
			if err != nil {
				fmt.Printf("❌ JSON Unmarshal error: %v\n", err)
				continue
			}

			fmt.Printf("✅ EventConsumer Description: %v\n", data["description"])
			fmt.Printf("✅ EventConsumer DueDate: %v\n", data["dueDate"])
			fmt.Printf("✅ EventConsumer id: %v\n", data["id"])

			event.Ack()
		}

		fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("✅ All events processed")
	})

	fmt.Println("🛑 Event consumer stopped")
}
