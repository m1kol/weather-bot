package provider

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testCity = "Долгопрудный"
	unknownCity = "asdqwe"
	testApiKey = "testKey"
)

func TestProvider_GetForecast(t *testing.T) {
	p := NewProvider(testApiKey)
	_, err := p.GetForecast(unknownCity)
	if err == nil {
		t.Log("Error should occur, but got: ", err)
		t.Fail()
	}
	_, err = p.GetForecast(testCity)
	if err != nil {
		t.Log("Error should be nil, but got:", err)
		t.Fail()
	}
}

func handleRequest(t *testing.T, city string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		require.Equal(t, testApiKey, req.URL.Query().Get("appid"), "Invalid authentication key")
		require.Equal(t, testCity, req.URL.Query().Get("q"))
	}
}
