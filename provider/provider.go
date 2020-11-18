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

type Provider struct {
	httpClient *resty.Client
	config map[string]string
}

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

func (prov *Provider) GetForecast(city string) (Forecast, error) {
	prov.config["q"] = city
	resp, err := prov.httpClient.R().SetQueryParams(prov.config).Get(weatherUrl)
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
	prov.config["lon"] = strconv.FormatFloat(cityInfo.Coord.Lon, 'f', -1, 32)
	prov.config["lat"] = strconv.FormatFloat(cityInfo.Coord.Lat, 'f', -1, 32)
	prov.config["exclude"] = "minutely,hourly,alerts,current"
	delete(prov.config, "q")

	resp, err = prov.httpClient.R().SetQueryParams(prov.config).Get(onecallUrl)
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
	prov.config["q"] = city
	prov.config["cnt"] = strconv.Itoa(8*days)
	resp, err := prov.httpClient.R().SetQueryParams(prov.config).Get(forecastUrl)
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