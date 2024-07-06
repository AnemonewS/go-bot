package main

import (
	"context"
	"flag"
	"log"
	tgClient "telegram-go/client/telegram"
	event_consumer "telegram-go/consumer/event-consumer"
	"telegram-go/events/telegram"
	"telegram-go/storage/postgresql"
)

const (
	tgHost         = "api.telegram.org"
	postgresqlPath = "user=postgres dbname=pqgotest sslmode=disable"
	batchSize      = 100
)

// TODO: improve error handling
func main() {
	//s := files.New(storagePath)
	s, err := postgresql.New(postgresqlPath)
	if err != nil {
		log.Fatalf("can't connect to the database", err)
	}
	if err := s.InitDatabase(context.TODO()); err != nil {
		log.Fatalf("can't initialize the database", err)
	}

	processor := telegram.New(
		tgClient.New(tgHost, mustToken()),
		s,
	)
	log.Printf("Service started")
	consumer_ := event_consumer.New(processor, processor, batchSize)
	if err := consumer_.Start(); err != nil {
		log.Fatal("Service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
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
