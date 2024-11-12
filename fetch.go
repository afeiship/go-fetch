package fetch

// https://github.com/go-zoox/fetch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Headers map[string]string
type Query map[string]string
type Params map[string]string
type Body interface{}

type Config struct {
	Url     string
	Method  string
	Headers Headers
	Query   Query
	Params  Params
	Body    Body
}

func Get(baseURL string, config *Config) (string, error) {
	// Replace params in the URL template
	for key, value := range config.Params {
		placeholder := fmt.Sprintf("{%s}", key)
		baseURL = strings.ReplaceAll(baseURL, placeholder, value)
	}

	// Parse the base URL to work with queries
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	// Add query parameters
	q := u.Query()
	for key, value := range config.Query {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	// Create the request
	req, err := http.NewRequest("GET", u.String(), nil)
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

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// 关闭响应
	resp.Body.Close()

	return string(body), nil
}

func Post(config Config) (string, error) {
	// TODO: implement Post method
	return "", nil
}

func Upload(config Config, filePath string) (string, error) {
	// TODO: implement Uploader method
	return "", nil
}
