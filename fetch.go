package fetch

// https://github.com/go-zoox/fetch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/afeiship/go-reader"
	nx "github.com/afeiship/nx/lib"
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

	// for upload
	ReaderType       reader.FileType
	ReaderSource     string
	MultipartOptions *nx.MultipartOptions
}

func Get(baseURL string, config *Config) (string, error) {
	u, _ := buildURL(baseURL, config)

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

func Upload(baseURL string, config *Config) (string, error) {
	u, _ := buildURL(baseURL, config)

	opts1 := reader.Options{
		Type:   config.ReaderType,
		Source: config.ReaderSource,
	}

	fileReader, err := reader.NewReader(&opts1)

	if err != nil {
		return "", fmt.Errorf("error creating file reader: %w", err)
	}

	opts2 := nx.MultipartOptions{
		FileReader:    fileReader,
		FieldName:     config.MultipartOptions.FieldName,
		FileFieldName: config.MultipartOptions.FileFieldName,
		ExtraFields:   config.MultipartOptions.ExtraFields,
	}

	multipartBody, err := nx.CreateMultipartRequestBody(&opts2)

	if err != nil {
		return "", fmt.Errorf("error creating multipart request body: %w", err)
	}

	req, err := http.NewRequest("POST", u.String(), multipartBody.Body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// set headers
	req.Header.Set("Content-Type", multipartBody.ContentType)
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	return string(respBody), nil
}

func buildURL(baseURL string, config *Config) (*url.URL, error) {
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

	return u, nil
}
