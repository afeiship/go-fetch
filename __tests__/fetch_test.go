package fetch

import (
	"github.com/tidwall/gjson"
	"testing"

	"github.com/afeiship/go-fetch"
)

func TestGet(f *testing.T) {
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
		f.Error(err)
	}

	resu := gjson.Get(res, "url")

	// check url + query
	if resu.String() != "https://www.httpbin.org/get?query1=value1&query2=value2" {
		f.Error("url is not correct", resu.String())
	}
	// check headers
	if gjson.Get(res, "headers.X-Custom-Header").String() != "aric" {
		f.Error("headers is not correct", gjson.Get(res, "headers.X-Custom-Header").String())
	}
}
