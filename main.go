package telegram_go

import (
	"flag"
	"log"
)

const (
	telegramBotHost = "api.telegram.org"
)

func main() {
	// token = flags.Get(token)
	//tgClient := telegram.New(telegramBotHost, mustToken())

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
