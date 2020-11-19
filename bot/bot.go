package bot

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/m1kol/weather-bot/provider"
	"log"
	"strings"
	"time"
)

const (
	forecastTemplate = `Дата и время: %v
Температура (ощущается):
* утром %.f (%.f) С 
* днём  %.f (%.f) С
* вечером  %.f (%.f) С
* ночью %.f (%.f) С
Погода: %v
Скорость ветра: %.1f м/с

`
	helpText = `I understand commands: /forecast.
Usage:
* /forecast city_name`
	welcomeText = `Hello, I'm a weather bot! I understand commands /forecast!
* /forecast command will give you a weather forecast for the day.`
)

type Bot struct {
	api *tgbotapi.BotAPI
	provider *provider.Provider
}

func NewBot(token string, provider *provider.Provider) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new bot: %w", err)
	}

	return &Bot{
		api: api,
		provider: provider,
	}, nil
}

func (bot *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.api.GetUpdatesChan(u)
	if err != nil {
		log.Printf("error getting updates: %v", err)
	}

	bot.ProcessMessage(updates)
}

func (bot *Bot) ProcessMessage(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Ignore non-command Messages
		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "start":
			msg.Text = welcomeText

		case "help":
			msg.Text = helpText

		case "forecast":
			city := update.Message.CommandArguments()
			if len(city) == 0 {
				msg.Text = helpText
				break
			}

			forecast, err := bot.provider.GetForecast(city)
			if err != nil {
				log.Printf("error getting a response from weather provider: %v", err)
			}

			builder := strings.Builder{}
			fmt.Fprintf(&builder, "Погода в городе %v\n\n", forecast.City)
			for _, day := range forecast.Daily {
				fmt.Fprintf(&builder, forecastTemplate,
					time.Unix(day.Dt, 0),
					day.Temp.Morn, day.FeelsLike.Morn,
					day.Temp.Day, day.FeelsLike.Day,
					day.Temp.Eve, day.FeelsLike.Eve,
					day.Temp.Night, day.FeelsLike.Night,
					day.Weather[0].Description, day.WindSpeed,
				)
			}

			msg.Text = builder.String()

		default:
			msg.Text = "I don't know that command.\n\n" + helpText
		}

		if _, err := bot.api.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}