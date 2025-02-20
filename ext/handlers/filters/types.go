package filters

import "github.com/anonyindian/gotgbot/v2"

type (
	CallbackQuery      func(cq *gotgbot.CallbackQuery) bool
	ChatMember         func(u *gotgbot.ChatMemberUpdated) bool
	ChosenInlineResult func(cir *gotgbot.ChosenInlineResult) bool
	InlineQuery        func(iq *gotgbot.InlineQuery) bool
	Message            func(msg *gotgbot.Message) bool
	Poll               func(poll *gotgbot.Poll) bool
	PollAnswer         func(pa *gotgbot.PollAnswer) bool
	ChatJoinRequest    func(cjr *gotgbot.ChatJoinRequest) bool
)
