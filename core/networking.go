package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const ApiCallFailed = "API call failed: "

func HttpRequest(method string, url string, headers map[string]string, data *strings.Reader) ([]byte, error) {
	// create Req instance
	req, err := http.NewRequest(method, url, data)
	if ErrorChecking(err) {
		return nil, CreateError(http.StatusInternalServerError, err.Error(), "")
	}

	// Add headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Call the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, CreateError(http.StatusInternalServerError, err.Error(), "")
	}

	// Read the body
	body, err := io.ReadAll(response.Body)
	if ErrorChecking(err) {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		bodyContent := struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}{}
		_ = json.Unmarshal(body, &bodyContent)
		message := ApiCallFailed + url + " - " + bodyContent.ErrorDescription
		fmt.Println(bodyContent)
		return nil, CreateError(response.StatusCode, message, bodyContent.Error)
	}

	defer func() { _ = response.Body.Close() }()
	return body, nil
}
