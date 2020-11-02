package bot

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot(apiKey string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to create a new bot: %w", err)
	}

	return &Bot{
		api: api,
	}, nil
}
