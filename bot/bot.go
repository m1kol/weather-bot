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
	helpText = `Я понимаю команды: /forecast.
Использование:
/forecast city_name`
	welcomeText = `Привет! Я бот прогноза погоды! Пока я понимаю только команду /forecast! С её помощью ты сможешь получить прогноз погоды на ближайшие дни.`
)

type Bot struct {
	api 			*tgbotapi.BotAPI
	provider 		provider.Provider
	subscriptions  	map[int64][]string
}

func NewBot(token string, provider provider.Provider) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new bot: %w", err)
	}

	return &Bot{
		api: api,
		provider: provider,
		subscriptions: map[int64][]string{},
	}, nil
}

func (bot *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.api.GetUpdatesChan(u)
	if err != nil {
		log.Printf("error getting updates: %v", err)
	}

	bot.processMessage(updates)
}

func (bot *Bot) processMessage(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Ignore non-command Messages
		if !update.Message.IsCommand() {
			continue
		}

		var msg tgbotapi.MessageConfig

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "start":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, welcomeText)

		case "help":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, helpText)

		case "forecast":
			msg = bot.processForecast(update)

		case "subscribe":
			msg = bot.processSubscribe(update)

		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Я не знаю такой команды.\n\n" + helpText)
		}

		if _, err := bot.api.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func (bot *Bot) processForecast(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	city := update.Message.CommandArguments()
	if len(city) == 0 {
		msg.Text = helpText
		return msg
	}

	forecast, err := bot.provider.GetForecast(city)
	if err != nil {
		log.Printf("error getting a response from weather provider: %v", err)
		msg.Text = fmt.Sprintf("Не удалось получить погоду в городе %v", city)
		return msg
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

	return msg
}

func (bot *Bot) processSubscribe(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	city := update.Message.CommandArguments()
	if len(city) == 0 {
		msg.Text = helpText
		return msg
	}

	_, err := bot.provider.GetForecast(city)
	if err != nil {
		msg.Text = fmt.Sprintf("Не удаётся получить погоду в городе %v, проверьте правильность написания.", city)
		return msg
	}

	bot.subscriptions[update.Message.Chat.ID] = append(bot.subscriptions[update.Message.Chat.ID], city)
	msg.Text = "Вы были успешно подписаны на прогноз погоды! Вы будете получать прогноз каждый день в 9 утра."

	return msg
}
