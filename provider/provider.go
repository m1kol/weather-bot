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
	weatherUrl = "https://api.openweathermap.org/data/2.5/weather"
	onecallUrl = "https://api.openweathermap.org/data/2.5/onecall"
)

//go:generate mockgen -destination mock/provider.go -package mock . Provider_
type Provider_ interface {
	GetForecast(city string) (Forecast, error)
}

type Provider struct {
	httpClient *resty.Client
}

func NewProvider(apiKey string) *Provider {
	client := resty.New()
	client.SetQueryParams(map[string]string{
		"lang": "ru",
		"units": "metric",
		"appid": apiKey,
	})

	return &Provider{httpClient: client}
}

func (prov *Provider) GetForecast(city string) (Forecast, error) {
	resp, err := prov.httpClient.R().SetQueryParam("q", city).Get(weatherUrl)
	if err != nil {
		return Forecast{}, fmt.Errorf("could not perform GET request: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return Forecast{}, fmt.Errorf("unexpected status code: %v", resp.StatusCode())
	}

	var cityInfo City
	if err := json.Unmarshal(resp.Body(), &cityInfo); err != nil {
		return Forecast{}, fmt.Errorf("failed to unmarshall data: %w", err)
	}

	resp, err = prov.httpClient.R().SetQueryParams(map[string]string{
		"lon": strconv.FormatFloat(cityInfo.Coord.Lon, 'f', -1, 32),
		"lat": strconv.FormatFloat(cityInfo.Coord.Lat, 'f', -1, 32),
		"exclude": "minutely,hourly,alerts,current",
	}).Get(onecallUrl)
	if err != nil {
		return Forecast{}, fmt.Errorf("could not perform GET request: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return Forecast{}, fmt.Errorf("unexpected status code: %v", resp.StatusCode())
	}

	var ret Forecast
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return Forecast{}, fmt.Errorf("failed to unmarshall data: %v", err)
	}
	ret.City = city

	return ret, nil
}

func (prov *Provider) GetWeather(city string, days int) (Response, error) {
	resp, err := prov.httpClient.R().SetQueryParams(map[string]string{
		"q": city,
		"cnt": strconv.Itoa(8*days),
	}).Get(forecastUrl)
	if err != nil {
		return Response{}, fmt.Errorf("could not perform GET request: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return Response{}, fmt.Errorf("unexpected status code: %v", resp.StatusCode())
	}

	var ret Response
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return Response{}, fmt.Errorf("failed to unmarshall data: %w", err)
	}

	return ret, nil
}

type City struct {
	Name  string	`json:"name"`
	Coord struct {
		Lon	float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
}

type Forecast struct {
	City 	string
	Daily 	[]Daily
}

type Daily struct {
	Dt   int64 `json:"dt"`
	Temp struct {
		Day		float32 `json:"day"`
		Night	float32 `json:"night"`
		Eve		float32 `json:"eve"`
		Morn 	float32 `json:"morn"`
	} `json:"temp"`
	FeelsLike struct {
		Day		float32 `json:"day"`
		Night	float32 `json:"night"`
		Eve		float32 `json:"eve"`
		Morn 	float32 `json:"morn"`
	} `json:"feels_like"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	WindSpeed float32 `json:"wind_speed"`
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
	Pop	float32 `json:"pop"`
}
