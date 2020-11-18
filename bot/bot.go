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
	replyTemplate = `Дата и время: %v
Температура: %.1f С	Максимальная: %.1f С Минимальная: %.1f С Ощущается как: %.1f С
Погода: %v
Скорость ветра: %.1f м/с

`
	forecastTemplate = `Дата и время: %v
Температура: утром %.f С днём  %.f С вечером  %.f С ночью %.f С
Ощущается как: утром %.f С днём  %.f С вечером  %.f С ночью %.f С
Погода: %v
Скорость ветра: %.1f м/с

`
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
		fmt.Errorf("error getting updates: %v", err)
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

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "start":
			msg.Text = "Hello, I'm a weather bot! I understand commands /sayhi, /status and /weather! " +
				"/weather commands will give you a weather forecast for the day."

		case "help":
			msg.Text = "I understand /sayhi, /status and /weather."

		case "sayhi":
			msg.Text = "Hi :)"

		case "status":
			msg.Text = "I'm ok."

		case "weather":
			city := update.Message.CommandArguments()
			resp, err := bot.provider.GetWeather(city, 1)
			if err != nil {
				fmt.Errorf("error getting a response from weather provider: %v", err)
			}

			builder := &strings.Builder{}
			fmt.Fprintf(builder, "Weather in a city %v\n\n", resp.City.Name)
			for i := 0; i < len(resp.WeatherInfo); i++ {
				fmt.Fprintf(builder, replyTemplate,
					resp.WeatherInfo[i].Time,
					resp.WeatherInfo[i].Main.Temp, resp.WeatherInfo[i].Main.MaxTemp,
					resp.WeatherInfo[i].Main.MinTemp, resp.WeatherInfo[i].Main.FeelsLike,
					resp.WeatherInfo[i].Weather[0].Description, resp.WeatherInfo[i].Wind.Speed,
				)
			}
			msg.Text = builder.String()

		case "forecast":
			city := update.Message.CommandArguments()
			forecast, err := bot.provider.GetForecast(city)
			if err != nil {
				log.Printf("error getting a response from weather provider: %v", err)
			}

			builder := strings.Builder{}
			fmt.Fprintf(&builder, "Погода в городе %v\n\n", forecast.City)
			for _, day := range forecast.Daily {
				fmt.Fprintf(&builder, forecastTemplate,
					time.Unix(day.Dt, 0),
					day.Temp.Morn, day.Temp.Day, day.Temp.Eve, day.Temp.Night,
					day.FeelsLike.Morn, day.FeelsLike.Day, day.FeelsLike.Eve, day.FeelsLike.Night,
					day.Weather[0].Description, day.WindSpeed,
				)
			}
			msg.Text = builder.String()

		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.api.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}