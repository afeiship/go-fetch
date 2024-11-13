# go-fetch
> Go Fetch - A Powerful, Lightweight, Easy Http Client, inspired by Web Fetch API.

## installation
```sh
go get -u github.com/afeiship/go-fetch
```

## usage
- Get
- Post
- Upload

### Get
```go
package main

// GET
import (
	"fmt"
	"log"
	"github.com/tidwall/gjson"
    "github.com/afeiship/go-fetch"
)

func main() {
	res, err := fetch.Get("https://www.httpbin.org/get", &fetch.Config{
		Headers: map[string]string{
			"X-Custom-Header": "aric",
		},
		Params: map[string]string{
			"param1": "value1",
		},
		Query: map[string]string{
			"query1": "value1",
			"query2": "value2",
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	resu := gjson.Get(res, "url")
	fmt.Println(resu.String())
}
```