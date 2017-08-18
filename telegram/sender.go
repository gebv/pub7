package telegram

import (
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

type TelegramAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

func TelegramSenderWithTimeout(
	api TelegramAPI,
	timeout time.Duration,
) *telegramSenderWithTimeout {
	return &telegramSenderWithTimeout{
		api,
		timeout,
	}
}

type telegramSenderWithTimeout struct {
	TelegramAPI

	Timeout time.Duration
}

func (s *telegramSenderWithTimeout) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	time.Sleep(s.Timeout)
	return s.TelegramAPI.Send(c)
}
