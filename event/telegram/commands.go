package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"telegram-go/lib/e"
	"telegram-go/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text, username string, chatId int) (err error) {
	text = strings.TrimSpace(text)

	log.Printf("Got new command '%s' from %s", text, username)

	if isAddCmd(text) {
		return p.savePage(chatId, text, username)
	}
	switch text {
	case StartCmd:
		return p.sendHello(chatId)
	case RndCmd:
		return p.sendRandom(chatId, username)
	case HelpCmd:
		return p.sendHelp(chatId)
	default:
		return p.tgClient.SendMessage(chatId, UnknownMessage)
	}
}

func (p *Processor) savePage(chatId int, pageUrl, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }()
	//send := NewMessageSender(chatId, p.tgClient)

	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}
	isExists, err := p.storage.Exists(page)
	if err != nil {
		return err
	}
	if isExists {
		//return send(AlreadySavedMessage)
		return p.tgClient.SendMessage(chatId, AlreadySavedMessage)
	}
	if err := p.storage.Save(page); err != nil {
		return err
	}
	if err := p.tgClient.SendMessage(chatId, SuccessfullySaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendRandom(chatId int, username string) (err error) {
	defer func() { err = e.WrapIfErr("Can't send random", err) }()

	page, err := p.storage.ChoseRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tgClient.SendMessage(chatId, NoSavedMessage)
	}
	if err := p.tgClient.SendMessage(chatId, page.URL); err != nil {
		return err
	}
	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatId int) (err error) {
	return p.tgClient.SendMessage(chatId, messageHelp)
}

func (p *Processor) sendHello(chatId int) (err error) {
	return p.tgClient.SendMessage(chatId, messageHello)
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}

//func NewMessageSender(chatId int, c *telegram.Client) func(string) error {
//	return func(message string) error {
//		return c.SendMessage(chatId, message)
//	}
//}
