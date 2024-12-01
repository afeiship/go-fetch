package fetch

import (
	"fmt"
	"io"
	"net/http"
)

func Put(baseURL string, config *Config) (string, error) {
	u, _ := buildURL(baseURL, config)

	// set defaults
	setDefaults(config)

	// 构建请求体
	body, contentType, err := buildBody(config)
	if err != nil {
		return "", err
	}

	// 创建请求
	req, err := http.NewRequest("PUT", u.String(), body)
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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	// 关闭响应
	err = resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
