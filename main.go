package telegram_go

import (
	"flag"
	"log"
	tgClient "telegram-go/client/telegram"
	event_consumer "telegram-go/consumer/event-consumer"
	"telegram-go/events/telegram"
	"telegram-go/storage/files"
)

const (
	tgHost      = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	processor := telegram.New(
		tgClient.New(tgHost, mustToken()),
		files.New(storagePath),
	)
	log.Printf("Service started")
	consumer_ := event_consumer.New(processor, processor, batchSize)
	if err := consumer_.Start(); err != nil {
		log.Fatal("Service is stopped", err)
	}
	// client = telegram.New(host, token) - Общается с АПИ ТГ
	// fetcher = fetcher.New() - Обращается к АПИ ТГ и получает сообщения
	// processor = processor.New() - Обрабатывает полученные сообщения и отдаем нам ссылки
	// consumer.Start(fetcher, processor) - Получает и обрабатывает события

}

func mustToken() string {
	token := flag.String(
		"token-bot-token",
		"",
		"Telegram API access token",
	)
	flag.Parse()
	token_ := *token

	if token_ == "" {
		log.Fatal("Token parameter is required")
	}
	return token_
}
