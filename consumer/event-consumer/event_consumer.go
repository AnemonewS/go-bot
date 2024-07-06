package event_consumer

import (
	"log"
	"math"
	"telegram-go/events"
	"telegram-go/lib/e"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		updateEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("Consumer fetch data error: %s", err)
			continue
		}
		if len(updateEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
	}
}

// Потеря событий: фоллбэк, ретрай
// Обработка пачки, получение ошибок: параллельная обработка, завершать после первой ошибки, счетчик ошибок
func (c *Consumer) handleEvents(events []events.Event) error {
	var unsuccessfulEvents int
	totalEvents := len(events)
	maxHandledError := int(math.Ceil(float64(totalEvents) / 3.0)) // 1/3

	for _, ev := range events {
		log.Printf("Got new events: %s", ev.Text)
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
