package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
)

type HttpEntity struct {
	Body    interface{}
	Headers map[string]string
}

// Exchange is a method to exchange data between 3rd party APIs
func Exchange(url string, method string, entity HttpEntity, responseStructPtr interface{}) error {
	// TODO: Add logging
	// First we need to create a new http client
	client := &http.Client{}

	// Then we need to create a new request
	requestBody, err := json.Marshal(entity.Body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	// Then we need to add the headers
	for key, value := range entity.Headers {
		req.Header.Add(key, value)
	}

	// Then we need to execute the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Then we need to check if response is same type as expected
	if responseStructPtr != nil {
		// TODO: How to validate if response is same type as expected?

		// Validate if response can be converted to expected type
		err = json.NewDecoder(resp.Body).Decode(responseStructPtr)
		if err != nil {
			return err
		}

		return nil
	}

	responseJson := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&responseJson)
	if err != nil {
		return err
	}

	return nil
}
