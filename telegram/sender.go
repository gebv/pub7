package telegram

import (
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

type TelegramAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

func TelegramSenderWithTimeout(
	api *tgbotapi.BotAPI,
	timeout time.Duration,
) *telegramSenderWithTimeout {
	return &telegramSenderWithTimeout{
		api,
		timeout,
	}
}

type telegramSenderWithTimeout struct {
	*tgbotapi.BotAPI

	Timeout time.Duration
}

func (s *telegramSenderWithTimeout) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	time.Sleep(s.Timeout)
	return s.BotAPI.Send(c)
}
