package fetch

import (
	"fmt"
	"github.com/afeiship/go-reader"
	nx "github.com/afeiship/nx/lib"
	"io"
	"net/http"
)

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