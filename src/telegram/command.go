package telegram

import tb "gopkg.in/telebot.v3"

const (
	START_CMD = "/start"
)

var Cmds = []tb.Command{
	{
		Text:        START_CMD,
		Description: "开始",
	},
}
