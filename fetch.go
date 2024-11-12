package fetch

// https://github.com/go-zoox/fetch

import (
	"fmt"
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

func SayHi() {
	fmt.Println("Hi from go-fetch")
}

func Get(config Config) (string, error) {
	// TODO: implement Get method
	return "", nil
}

func Post(config Config) (string, error) {
	// TODO: implement Post method
	return "", nil
}

func Upload(config Config, filePath string) (string, error) {
	// TODO: implement Uploader method
	return "", nil
}
