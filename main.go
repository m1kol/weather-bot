package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/go-resty/resty"
)

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
	prov.config["city"] = city
	prov.config["cnt"] = strconv.Itoa(8*days)

	fmt.Println(prov.config)

	resp, err := prov.httpClient.R().SetQueryParams(prov.config).Get(forecastUrl)
	fmt.Println(resp)
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
	city 		string 			`json:"city.name"`
	weatherInfo []WeatherInfo 	`json:"list"`
}

type WeatherInfo struct {
	time        time.Time `json:"dt"`
	temp        float32   `json:"main.temp"`
	maxTemp     float32   `json:"main.temp_max"`
	minTemp     float32   `json:"main.temp_min"`
	feels_like  float32   `json:"main.feels_like"`
	weatherType string    `json:"weather.description"`
	windSpeed   float32   `json:"wind.speed"`
	pop         float32   `json:"pop"`
}

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

func main() {
	provider := NewProvider("myKey")
	city := "Долгопрудный"
	tmp, err := provider.GetWeather(city, 5)
	if err != nil {
		log.Fatalf("Failed to get weather information: %w", err)
	}

	fmt.Printf("Погода в %s\n", tmp.city)
	for i := 0; i < len(tmp.weatherInfo); i++ {
		fmt.Println(tmp.weatherInfo[i])
	}
}
