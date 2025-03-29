package fetch

import (
	"fmt"
	"io"
	"net/http"
)

func Request(method string, baseURL string, config *Config) (string, error) {
	u, _ := buildURL(baseURL, config)

	// set defaults
	setDefaults(config)

	// Build request body if not GET or DELETE
	var body io.Reader
	var contentType string
	var err error

	if method != "GET" && method != "DELETE" {
		body, contentType, err = buildBody(config)
		if err != nil {
			return "", err
		}
	}

	// Create the request
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// Set content type if present
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %v", err)
	}

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Close response
	err = resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}