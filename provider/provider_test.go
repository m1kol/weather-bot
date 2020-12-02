package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"net/http"
	//"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const (
	testCity = "myCity"
	testCityUnknown = "asdqwezxc"
	testApiKey = "testKey"
	testDataFile = "test_data.json"
)

func TestOWMProvider_GetForecast(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

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

	weatherQuery := map[string]string{
		"lang": "ru",
		"units": "metric",
		"appid": testApiKey,
		"q": testCity,
	}
	httpmock.RegisterResponderWithQuery("GET", weatherUrl,
		weatherQuery, httpmock.NewBytesResponder(200, []byte(`{"name":myCity,"coord":{"lon":-94.04,"lat":33.44}}`)))

	fmt.Println("THIS")

	var testForecast Forecast
	err = json.Unmarshal(testData, &testForecast)
	require.NoError(t, err, "Failed to unmarshall test data")
	testForecast.City = testCity
	fmt.Println("THAT")

	p := NewProvider(testApiKey)
	fmt.Println("LALA")
	forecast, err := p.GetForecast(testCity)
	fmt.Println(forecast, err)
	require.NoError(t, err, "Failed to get weather forecast")
	require.Equal(t, forecast.City, testForecast.City, "Incorrect city name")
	require.Equal(t, testForecast.Daily[0].Temp.Day, forecast.Daily[0].Temp.Day,
		"Incorrect temperature response",
	)
}

//func handleForecastRequest(t *testing.T) http.HandlerFunc {
//	return func(writer http.ResponseWriter, request *http.Request) {
//		require.Equal(t, testApiKey, request.URL.Query().Get("appid"), "Invalid API key")
//
//		data, err := ioutil.ReadFile(testDataFile)
//		require.NoError(t, err, "Error reading json file")
//		_, _ = writer.Write(data)
//	}
//}
