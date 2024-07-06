package event_consumer

import (
	"log"
	"math"
	"telegram-go/event"
	"telegram-go/lib/e"
	"time"
)

type Consumer struct {
	fetcher   event.Fetcher
	processor event.Processor
	batchSize int
}

func New(fetcher event.Fetcher, processor event.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		events, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("Consumer fetch data error: %s", err)
			continue
		}
		if len(events) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
	}
}

// Потеря событий: фоллбэк, ретрай
// Обработка пачки, получение ошибок: параллельная обработка, завершать после первой ошибки, счетчик ошибок
func (c *Consumer) handleEvents(events []event.Event) error {
	var unsuccessfulEvents int
	totalEvents := len(events)
	maxHandledError := int(math.Ceil(float64(totalEvents) / 3.0)) // 1/3

	for _, ev := range events {
		log.Printf("Got new event: %s", ev.Text)
		if err := c.processor.Process(ev); err != nil {
			log.Printf("can't handle error: %s", err.Error())
			unsuccessfulEvents += 1
			if unsuccessfulEvents == maxHandledError {
				return e.WrapError("Got max handled errors", err)
			}
			continue
		}
	}
	return nil
}
