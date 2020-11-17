package provider

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strconv"
)

const (
	forecastUrl = "https://api.openweathermap.org/data/2.5/forecast"
	//weatherUrl = "https://api.openweathermap.org/data/2.5/weather"
)

type Provider struct {
	httpClient *resty.Client
	config map[string]string
}

//type Provider interface {
//	GetWeather(city string, days int)
//}

func NewProvider(apiKey string) *Provider {
	return &Provider{
		httpClient: resty.New(),
		config: map[string]string{
			"lang": "ru",
			"units": "metric",
			"appid": apiKey,
		},
	}
}

func (prov *Provider) GetWeather(city string, days int) (Response, error) {
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
