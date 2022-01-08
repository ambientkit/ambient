// Package requestclient provides HTTP request sending using JSON structs.
package requestclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// RequestClient is for making HTTP requests.
type RequestClient struct {
	baseURL     string
	bearerToken string
}

// New returns a request client.
func New(baseURL string, bearerToken string) *RequestClient {
	return &RequestClient{
		baseURL:     baseURL,
		bearerToken: bearerToken,
	}
}

// Get makes a GET request.
func (c *RequestClient) Get(urlSuffix string, returnData interface{}) error {
	// Ensure supported returnData was passed in (should be pointer).
	if returnData != nil {
		v := reflect.ValueOf(returnData)
		if v.Kind() != reflect.Ptr {
			return errors.New("data must pass a pointer, not a value")
		}
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v%v", c.baseURL, urlSuffix), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err.Error())
	}

	if len(c.bearerToken) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.bearerToken))
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error creating client: %v", err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %v", err.Error())
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error GET response (%v): %v", resp.StatusCode, err.Error())
	}

	if returnData != nil {
		err = json.Unmarshal(body, returnData)
		if err != nil {
			return fmt.Errorf("error unmarshal: %v", err.Error())
		}
	}

	return nil
}

// Post makes a POST request.
func (c *RequestClient) Post(urlSuffix string, sendData interface{}, returnData interface{}) error {
	// Ensure supported returnData was passed in (should be pointer).
	if returnData != nil {
		v := reflect.ValueOf(returnData)
		if v.Kind() != reflect.Ptr {
			return errors.New("data must pass a pointer, not a value")
		}
	}

	var err error
	var req *http.Request
	if sendData != nil {
		// Send data with the request if passed in.
		sendJSON, err := json.Marshal(sendData)
		if err != nil {
			return err
		}

		req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%v%v", c.baseURL, urlSuffix), bytes.NewReader(sendJSON))
		if err != nil {
			return fmt.Errorf("error creating request: %v", err.Error())
		}
	} else {
		// Don't send data in with the request if passed in.
		req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%v%v", c.baseURL, urlSuffix), nil)
		if err != nil {
			return fmt.Errorf("error creating request: %v", err.Error())
		}
	}

	if len(c.bearerToken) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.bearerToken))
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error creating client: %v", err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %v", err.Error())
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error GET response (%v): %v", resp.StatusCode, err.Error())
	}

	if returnData != nil {
		err = json.Unmarshal(body, returnData)
		if err != nil {
			return fmt.Errorf("error unmarshal: %v", err.Error())
		}
	}

	return nil
}
