package telegram

import (
	"errors"
	"telegram-go/client/telegram"
	"telegram-go/event"
	"telegram-go/lib/e"
	"telegram-go/storage"
)

var (
	UnknownEventTypeError = errors.New("unknown event type")
	UnknownMetaTypeError  = errors.New("unknown meta type")
)

type Processor struct {
	tgClient *telegram.Client
	offset   int
	// storage -> to save links
	storage storage.Storage
}

type Meta struct {
	ChatId   int
	Username string
}

func New(client *telegram.Client, s storage.Storage) *Processor {
	return &Processor{
		tgClient: client,
		storage:  s,
	}
}

func (p *Processor) Fetch(limit int) ([]event.Event, error) {
	updates, err := p.tgClient.Updates(p.offset, limit)
	if err != nil {
		return nil, e.WrapError("can't fetch updates", err)
	}
	if len(updates) == 0 {
		return nil, nil
	}
	res := make([]event.Event, 0, len(updates))

	for _, update := range updates {
		res = append(res, makeEvent(update))
	}
	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func makeEvent(update telegram.Update) event.Event {
	updateType := fetchType(update)

	res := event.Event{
		Type: fetchType(update),
		Text: fetchText(update),
	}
	if updateType == event.Message {
		res.Meta = Meta{
			update.Message.Chat.ID,
			update.Message.From.Username,
		}
	}
	return res
}

func fetchType(update telegram.Update) event.Type {
	if update.Message == nil {
		return event.Unknown
	}
	return event.Message
}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}

func (p *Processor) Process(ev event.Event) error {
	switch ev.Type {
	case event.Message:
		return p.ProcessMessage(ev)
	default:
		return e.WrapError("can't process message", UnknownEventTypeError)
	}

}

func (p *Processor) ProcessMessage(ev event.Event) error {
	meta, err := prepareMeta(ev)

	if err != nil {
		return e.WrapError("can't process message", err)
	}
	if err := p.doCmd(ev.Text, meta.Username, meta.ChatId); err != nil {
		return e.WrapError("can't process message", err)
	}
	return nil

}

func prepareMeta(ev event.Event) (Meta, error) {
	meta, ok := ev.Meta.(Meta)

	if !ok {
		return Meta{}, e.WrapError("can't get meta", UnknownMetaTypeError)
	}
	return meta, nil
}
