package gotgbot

import (
	"time"

	"github.com/anonyindian/gotgbot/v2/request"
)

//go:generate go run ./scripts/generate

// Bot is the core Bot object used to send and receive messages.
type Bot struct {
	User
	Request request.RequestType
}

// BotOpts declares all optional parameters for the NewBot function.
type BotOpts struct {
	APIURL      string
	Request     request.RequestType
	GetTimeout  time.Duration
	PostTimeout time.Duration
}

// NewBot returns a new Bot struct populated with the necessary defaults.
func NewBot(token string, opts *BotOpts) (*Bot, error) {
	b := Bot{}

	getTimeout := request.DefaultGetTimeout
	postTimeout := request.DefaultPostTimeout
	apiUrl := request.DefaultAPIURL
	if opts != nil {
		b.Request = opts.Request
		apiUrl = opts.APIURL
		getTimeout = opts.GetTimeout
		postTimeout = opts.PostTimeout
	}
	b.Request = &request.BuiltinHttp{
		ApiUrl:      apiUrl,
		BotToken:    token,
		PostTimeout: postTimeout,
		GetTimeout:  getTimeout,
	}
	u, err := b.GetMe()
	if err != nil {
		return nil, err
	}

	b.User = *u
	return &b, nil
}
