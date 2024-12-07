package fetch

// https://github.com/go-zoox/fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/afeiship/go-reader"
	nx "github.com/afeiship/nx/lib"
)

type Headers map[string]string
type Query map[string]string
type Params map[string]string
type Body interface{}
type DataType string

type Config struct {
	Url      string
	Method   string
	DataType DataType
	Headers  Headers
	Query    Query
	Params   Params
	Body     Body

	// for upload
	ReaderType       reader.FileType
	ReaderSource     string
	MultipartOptions *nx.MultipartOptions
}

const (
	URLENCODED   DataType = "urlencoded"
	JSON         DataType = "json"
	OCTET_STREAM DataType = "octet-stream"
)

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
	case JSON:
		jsonData, err := json.Marshal(config.Body)
		if err != nil {
			return nil, "", err
		}
		body = bytes.NewBuffer(jsonData)
		contentType = "application/json"

	case URLENCODED:
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

	case OCTET_STREAM:
		body = strings.NewReader(config.Body.(string))
		contentType = "application/octet-stream"
		return body, contentType, nil

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
		config.DataType = URLENCODED
	}
}
