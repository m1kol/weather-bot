package main

import (
	"flag"
	"fmt"
	"github.com/m1kol/weather-bot/bot"
	"github.com/m1kol/weather-bot/provider"
)

// Initialization and main

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

	//log.Printf("euthorized on account %v", bot.api.Self.UserName)

	bot.Run()

	//city := "Долгопрудный"
	//
	//res, err := provider.GetWeather(city, 5)
	//if err != nil {
	//	log.Fatalf("Failed to get weather information: %w", err)
	//}
	//
	//fmt.Println(res.City.Name)
	//fmt.Println(res)
}
