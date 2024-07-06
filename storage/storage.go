package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"telegram-go/lib/e"
)

type Storage interface {
	Save(page *Page) error
	ChoseRandom(userName string) (*Page, error)
	Remove(page *Page) error
	Exists(page *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("there are no saved pages")

type Page struct {
	URL      string
	UserName string
}

func (page *Page) Hash() (string, error) {
	h := sha1.New()
	const calculateHashMsg = "Can't calculate hash"

	if _, err := io.WriteString(h, page.URL); err != nil {
		return "", e.WrapError(calculateHashMsg, err)
	}
	if _, err := io.WriteString(h, page.UserName); err != nil {
		return "", e.WrapError(calculateHashMsg, err)
	}
	return fmt.Sprint("%x", h.Sum(nil)), nil
}
