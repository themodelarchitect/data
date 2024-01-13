package client

import (
	"bytes"
	"net/url"
	"testing"
)

// docker run -p 80:80 kennethreitz/httpbin

func TestGetHTTPS(t *testing.T) {
	_, err := Get("https://devstreaming-cdn.apple.com/videos/streaming/examples/img_bipbop_adv_example_ts/master.m3u8", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHead(t *testing.T) {
	resp, err := Head("https://golang.org")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp.StatusCode)
	t.Log(resp.ContentLength)

}

func TestGet(t *testing.T) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	// the query parameters to pass
	queryParameters := url.Values{}
	queryParameters.Add("foo", "bar")
	res, err := Get("http://127.0.0.1/get", headers, queryParameters)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(res))
}

func TestPost(t *testing.T) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	// the body to pass
	body := bytes.NewBufferString(`{"name": "test"}`)
	res, err := Post("http://127.0.0.1/post", headers, body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(res))
}

func TestPut(t *testing.T) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	// the body to pass
	body := bytes.NewBufferString(`{"name": "test"}`)
	res, err := Put("http://127.0.0.1/put", headers, body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(res))
}

func TestPatch(t *testing.T) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	// the body to pass
	body := bytes.NewBufferString(`{"name": "test"}`)
	res, err := Patch("http://127.0.0.1/patch", headers, body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(res))
}

func TestDelete(t *testing.T) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	res, err := Delete("http://127.0.0.1/delete", headers)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(res))
}
