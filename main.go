package main

import (
	"flag"
	"fmt"
	"github.com/m1kol/weather-bot/bot"
	"github.com/m1kol/weather-bot/provider"
)

var (
	apiKey		string
	botToken	string
)

func init() {
	flag.StringVar(&apiKey, "api-key", "", "OpenWeather API token")
	flag.StringVar(&botToken, "bot-token", "", "Telegram Bot API token")
	flag.Parse()
}

func main() {
	provider := provider.NewProvider(apiKey)
	bot, err := bot.NewBot(botToken, provider)
	if err != nil {
		fmt.Errorf("error creating bot: %v", err)
	}

	bot.Run()
}
