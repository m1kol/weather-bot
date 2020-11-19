package main

import (
	"flag"
	"github.com/m1kol/weather-bot/bot"
	"github.com/m1kol/weather-bot/provider"
	"log"
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
	prov := provider.NewProvider(apiKey)
	b, err := bot.NewBot(botToken, prov)
	if err != nil {
		log.Printf("error creating bot: %v", err)
	}

	b.Run()
}
