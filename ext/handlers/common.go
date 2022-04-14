package handlers

import (
	"github.com/anonyindian/gotgbot/v2"
	"github.com/anonyindian/gotgbot/v2/ext"
)

type Response func(b *gotgbot.Bot, ctx *ext.Context) error
