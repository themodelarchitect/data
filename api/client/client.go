package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type requestMethod byte

const (
	GET requestMethod = iota
	POST
	PUT
	PATCH
	DELETE
)

func (method requestMethod) String() string {
	switch method {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case PATCH:
		return "PATCH"
	case DELETE:
		return "DELETE"
	default:
		return fmt.Sprintf("Unknown(%d)", method)
	}
}

func Get(url string, headers map[string]string, queryParams url.Values) ([]byte, error) {
	return request(GET, url, headers, queryParams, nil)
}

func Post(url string, headers map[string]string, body io.Reader) ([]byte, error) {
	return request(POST, url, headers, nil, body)
}

func Put(url string, headers map[string]string, body io.Reader) ([]byte, error) {
	return request(PUT, url, headers, nil, body)
}

func Patch(url string, headers map[string]string, body io.Reader) ([]byte, error) {
	return request(PATCH, url, headers, nil, body)
}

func Delete(url string, headers map[string]string) ([]byte, error) {
	return request(DELETE, url, headers, nil, nil)
}

func request(method requestMethod, fullUrl string, headers map[string]string, queryParams url.Values, body io.Reader) ([]byte, error) {
	var b []byte

	u, err := url.Parse(fullUrl)
	if err != nil {
		return b, err
	}

	// if it's a GET, we need to append the query parameters.
	if method == GET {
		q := u.Query()

		for k, v := range queryParams {
			// this depends on the type of api, you may need to do it for each of v
			q.Set(k, strings.Join(v, ","))
		}
		// set the query to the encoded parameters
		u.RawQuery = q.Encode()
	}

	// create a request
	req, err := http.NewRequest(method.String(), u.String(), body)
	if err != nil {
		return b, err
	}

	// NOTE this !! -You need to set Req.Close to true (the defer on resp.Body.Close() syntax used in the examples is not enough)
	req.Close = true

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	// send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return b, err
	}

	defer func(Body io.ReadCloser) {
		// Ignore error explicitly
		_ = Body.Close()
	}(resp.Body)

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return b, err
	}

	return b, nil
}
