package client

import (
	"bytes"
	"log"
	"net/url"
	"testing"
)

// docker run -p 80:80 kennethreitz/httpbin

func TestGet(t *testing.T) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	// the query parameters to pass
	queryParameters := url.Values{}
	queryParameters.Add("foo", "bar")
	res, err := Get("http://127.0.0.1/get", headers, queryParameters)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	log.Println(string(res))
}

func TestPost(t *testing.T) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	// the body to pass
	body := bytes.NewBufferString(`{"name": "test"}`)
	res, err := Post("http://127.0.0.1/post", headers, body)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	log.Println(string(res))
}

func TestPut(t *testing.T) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	// the body to pass
	body := bytes.NewBufferString(`{"name": "test"}`)
	res, err := Put("http://127.0.0.1/put", headers, body)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	log.Println(string(res))
}

func TestPatch(t *testing.T) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	// the body to pass
	body := bytes.NewBufferString(`{"name": "test"}`)
	res, err := Patch("http://127.0.0.1/patch", headers, body)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	log.Println(string(res))
}

func TestDelete(t *testing.T) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	res, err := Delete("http://127.0.0.1/delete", headers)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	log.Println(string(res))
}
