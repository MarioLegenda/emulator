package httpClient

import (
	"bytes"
	"net/http"
	"time"
)

type header struct {
	Name  string
	Value string
}

type request struct {
	Headers []header
	Url     string
	Method  string
	Body    []byte
}

type ClientOption func(f *http.Client)

func WithTransport(option http.RoundTripper) ClientOption {
	return func(client *http.Client) {
		client.Transport = option
	}
}

func WithTimeout(option time.Duration) ClientOption {
	return func(client *http.Client) {
		client.Timeout = option
	}
}

func WithJar(option http.CookieJar) ClientOption {
	return func(client *http.Client) {
		client.Jar = option
	}
}

func WithCheckRedirect(option func(req *http.Request, via []*http.Request) error) ClientOption {
	return func(client *http.Client) {
		client.CheckRedirect = option
	}
}

func NewHeader(name string, value string) header {
	return header{
		Name:  name,
		Value: value,
	}
}

func NewClient(opts ...ClientOption) *http.Client {
	client := &http.Client{}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func NewRequest(request request) (*http.Request, error) {
	r, err := http.NewRequest(request.Method, request.Url, bytes.NewBuffer(request.Body))

	if err != nil {
		return nil, err
	}

	if len(request.Headers) != 0 {
		for _, v := range request.Headers {
			r.Header.Set(v.Name, v.Value)
		}
	}

	return r, nil
}

func Make(request *http.Request, client *http.Client) (*http.Response, error) {
	return client.Do(request)
}
