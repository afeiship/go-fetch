package fetch

// https://github.com/go-zoox/fetch

import (
	"bytes"
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

func Get(baseURL string, config *Config) (*http.Response, error) {
	// Replace params in the URL template
	for key, value := range config.Params {
		placeholder := fmt.Sprintf("{%s}", key)
		baseURL = strings.ReplaceAll(baseURL, placeholder, value)
	}

	// Parse the base URL to work with queries
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	// Add query parameters
	q := u.Query()
	for key, value := range config.Query {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	// Create the request
	req, err := http.NewRequest(config.Method, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// Add body if specified and method is not GET
	if config.Body != nil && config.Method != "GET" {
		var body io.Reader
		switch b := config.Body.(type) {
		case string:
			body = strings.NewReader(b)
		case []byte:
			body = bytes.NewReader(b)
		case io.Reader:
			body = b
		default:
			return nil, fmt.Errorf("unsupported body type")
		}
		req.Body = io.NopCloser(body)
	}

	// Send the request
	client := &http.Client{}
	return client.Do(req)
}

func Post(config Config) (string, error) {
	// TODO: implement Post method
	return "", nil
}

func Upload(config Config, filePath string) (string, error) {
	// TODO: implement Uploader method
	return "", nil
}
