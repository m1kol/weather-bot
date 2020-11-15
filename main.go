package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Provider classes and methods

const (
	forecastUrl = "https://api.openweathermap.org/data/2.5/forecast"
	//weatherUrl = "https://api.openweathermap.org/data/2.5/weather"
)

type Provider_ struct {
	httpClient *resty.Client
	config map[string]string
}

type Provider interface {
	GetWeather(city string, days int)
}

func NewProvider(apiKey string) *Provider_ {
	return &Provider_{
		httpClient: resty.New(),
		config: map[string]string{
			"lang": "ru",
			"units": "metric",
			"appid": apiKey,
		},
	}
}

func (prov *Provider_) GetWeather(city string, days int) (Response, error) {
	prov.config["q"] = city
	prov.config["cnt"] = strconv.Itoa(8*days)

	resp, err := prov.httpClient.R().SetQueryParams(prov.config).Get(forecastUrl)
	if err != nil {
		return Response{}, fmt.Errorf("could not perform GET request: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return Response{}, fmt.Errorf("unexpected status code: %w", resp.StatusCode())
	}

	var ret Response
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return Response{}, fmt.Errorf("failed to unmarshall data: %w", err)
	}

	return ret, nil
}

type Response struct {
	City struct {
		Name string `json:"name"`
	} `json:"city"`
	WeatherInfo	[]WeatherInfo `json:"list"`
}

type WeatherInfo struct {
	Time string	`json:"dt_txt"`
	Main struct {
		Temp        float32 `json:"temp"`
		MaxTemp     float32 `json:"temp_max"`
		MinTemp     float32 `json:"temp_min"`
		FeelsLike	float32 `json:"feels_like"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float32 `json:"speed"`
	} `json:"wind"`
	Pop         float32 `json:"pop"`
}

// Bot classes and methods

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("Failed to create a new bot: %w", err)
	}

	return &Bot{
		api: api,
	}, nil
}

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
	provider := NewProvider(apiKey)
	city := "Долгопрудный"
	res, err := provider.GetWeather(city, 5)
	if err != nil {
		log.Fatalf("Failed to get weather information: %w", err)
	}

	fmt.Println(res.City.Name)
	fmt.Println(res)
}
