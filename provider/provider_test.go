package provider

import (
	"testing"
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