package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const (
	testCity = "myCity"
	testApiKey = "testKey"
	testDataFile = "testdata/test_data.json"
)

func TestOWMProvider_GetForecast(t *testing.T) {
	p := NewProvider(testApiKey)
	httpmock.ActivateNonDefault(p.httpClient.GetClient())
	defer httpmock.DeactivateAndReset()

	testCityInfo := City{
		Name: testCity,
		Coord: struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		}{Lon: -94.04, Lat: 33.44},
	}
	weatherQuery := map[string]string{
		"lang": "ru",
		"units": "metric",
		"appid": testApiKey,
		"q": testCity,
	}
	httpmock.RegisterResponderWithQuery("GET", weatherUrl,
		weatherQuery, httpmock.NewJsonResponderOrPanic(200, &testCityInfo))

	testData, err := ioutil.ReadFile(testDataFile)
	require.NoError(t, err, "Failed to read from a data file")
	onecallQuery := map[string]string{
		"lang": "ru",
		"units": "metric",
		"appid": testApiKey,
		"lon": "-94.04",
		"lat": "33.44",
		"exclude": "minutely,hourly,alerts,current",
	}
	httpmock.RegisterResponderWithQuery("GET", onecallUrl,
		onecallQuery, httpmock.NewBytesResponder(200, testData))


	var testForecast Forecast
	err = json.Unmarshal(testData, &testForecast)
	require.NoError(t, err, "Failed to unmarshall test data")
	testForecast.City = testCity

	forecast, err := p.GetForecast(testCity)
	fmt.Println(forecast)
	require.NoError(t, err, "Failed to get weather forecast")
	require.Equal(t, testForecast.City, forecast.City, "Incorrect city name")
	require.Equal(t, testForecast.Daily[0].Temp.Day, forecast.Daily[0].Temp.Day,
		"Incorrect temperature response",
	)
}
