package users

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func NewTgBot() (*tgbotapi.BotAPI, error) {
	tgBot, err := tgbotapi.NewBotAPI("apiKeyExample") // апи ключ тг бота
	return tgBot, errors.Wrap(err, "creating new tg bot")
}
