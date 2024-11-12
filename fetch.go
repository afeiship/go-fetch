package fetch

// https://github.com/go-zoox/fetch

import (
	"bytes"
	"encoding/json"
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
	Url      string
	Method   string
	DataType string
	Headers  Headers
	Query    Query
	Params   Params
	Body     Body

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

func Post(baseURL string, config *Config) (string, error) {
	u, _ := buildURL(baseURL, config)

	// set defaults
	setDefaults(config)

	// 构建请求体
	body, contentType, err := buildBody(config)
	if err != nil {
		return "", err
	}

	// 创建请求
	req, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		return "", err
	}

	// 设置请求头
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}
	// 设置内容类型
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
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

// ----------------------------- private functions -----------------------------

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

func buildBody(config *Config) (io.Reader, string, error) {
	var body io.Reader
	var contentType string

	switch config.DataType {
	case "json":
		jsonData, err := json.Marshal(config.Body)
		if err != nil {
			return nil, "", err
		}
		body = bytes.NewBuffer(jsonData)
		contentType = "application/json"

	case "urlencode":
		values := url.Values{}
		if formData, ok := config.Body.(map[string]string); ok {
			for key, value := range formData {
				values.Set(key, value)
			}
		} else {
			return nil, "", fmt.Errorf("urlencode body must be map[string]string")
		}
		body = strings.NewReader(values.Encode())
		contentType = "application/x-www-form-urlencoded"

	default:
		return nil, "", fmt.Errorf("unsupported DataType: %s", config.DataType)
	}

	return body, contentType, nil
}

func setDefaults(config *Config) {
	// 设置默认的 Headers、Query 和 Params 为空 map
	if config.Headers == nil {
		config.Headers = Headers{}
	}
	if config.Query == nil {
		config.Query = Query{}
	}
	if config.Params == nil {
		config.Params = Params{}
	}
	// 设置默认的 DataType 为 "urlencode"
	if config.DataType == "" {
		config.DataType = "urlencode"
	}
}
