package rest

import (
	"testing"
)

const sampleApiUrl = "https://reqres.in/api/"

func TestExchangeGetCorrect(t *testing.T) {
	// Send GET request to sample API
	// Check if response is expected
	type expectedMap struct {
		Data map[string]interface{} `json:"data"`
	}
	responseStruct := new(expectedMap)
	err := Exchange(sampleApiUrl+"users/2", GET, HttpEntity{}, responseStruct)
	if err != nil {
		t.Error(err)
	}

	if responseStruct.Data["id"] == nil {
		t.Error("Response is not expected")
	}
}
