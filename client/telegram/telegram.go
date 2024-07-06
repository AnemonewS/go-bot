package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"telegram-go/lib/e"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string // Telegram API host
	basePath string // tg-bot.com/bot<token>
	client   http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset, limit int) (updates []Update, err error) {
	query := url.Values{}
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	// Make request
	data, err := c.doRequest(getUpdatesMethod, query)
	if err != nil {
		return nil, err
	}
	var res UpdateResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) SendMessage(chatId int, text string) error {
	query := url.Values{}
	query.Add("chat_id", strconv.Itoa(chatId))
	query.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, query)
	if err != nil {
		return e.WrapError("can't send message", err)
	}
	return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("doRequest func: can't do request", err) }()

	url_ := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method), // instead of -> c.basePath + method
	}
	request, err := http.NewRequest(
		http.MethodGet,
		url_.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = query.Encode()
	response, err := c.client.Do(request)

	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
