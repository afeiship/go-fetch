package fetch

import (
	"fmt"
	"io"
	"net/http"
)

func Delete(baseURL string, config *Config) (string, error) {
	u, _ := buildURL(baseURL, config)

	// Create the request
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %v", err)
	}

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Close response
	err = resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(body), nil
}