package fetch

import (
	"fmt"
	"os"
	"testing"

	nx "github.com/afeiship/nx/lib"
	"github.com/tidwall/gjson"

	"github.com/afeiship/go-fetch"
	"github.com/afeiship/go-reader"
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

func TestUpload(f *testing.T) {
	var cookie = fmt.Sprintf("SUB=%s", os.Getenv("WEIBO_TOKEN"))
	var res, err = fetch.Upload("https://picupload.weibo.com/interface/pic_upload.php", &fetch.Config{
		ReaderType:   reader.File,
		ReaderSource: "./01.jpg",
		MultipartOptions: &nx.MultipartOptions{
			FieldName:     "pic1",
			FileFieldName: "01.jpg",
		},
		Headers: map[string]string{
			"Cookie": cookie,
		},
	})

	if err != nil {
		f.Error(err)
	}
	fmt.Println(res)
}

func TestPost(t *testing.T) {
	res, err := fetch.Post("https://www.httpbin.org/post", &fetch.Config{
		DataType: "urlencode",
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
		Body: map[string]string{
			"name": "aric",
			"age":  "25",
		},
	})

	if err != nil {
		t.Error(err)
	}

	fmt.Println("result is: ", res)

	// resu := gjson.Get(res, "url")
	// // check url + query
	// if resu.String() != "https://www.httpbin.org/post?query1=value1&query2=value2" {
	// 	t.Error("url is not correct", resu.String())
	// }
	// // check headers
	// if gjson.Get(res, "headers.X-Custom-Header").String() != "aric" {
	// 	t.Error("headers is not correct", gjson.Get(res, "headers.X-Custom-Header").String())
	// }
	// // check params
	// if gjson.Get(res, "form.param1").String() != "value1" {
	// 	t.Error("params is not correct", gjson.Get(res, "form.param1").String())
	// }
	// // check body
	// if gjson.Get(res, "json.name").String() != "aric" {
	// 	t.Error("body is not correct", gjson.Get(res, "json.name").String())
	// }
}
