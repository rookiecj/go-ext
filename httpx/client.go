package httpx

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

const (
	DefaultTimeout = 10 * time.Second
)

type Client interface {
	Do(method string, url url.URL, options ...RequestOption) (res *http.Response, data []byte, err error)
}

type RequestOption func(req *http.Request)

type Request interface {
}

//type Response struct {
//	Header     map[string][]string
//	Status     string
//	StatusCode int
//	Data       []byte
//}

type baseClient struct {
	client             http.Client
	bodyParsers        map[reflect.Type]BodyParser
	disableCompression bool
}

var DefaultClient = NewBuilder().Build()

func (c *baseClient) Do(method string, url url.URL, options ...RequestOption) (resp *http.Response, data []byte, err error) {
	// new request
	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		err = fmt.Errorf("failed to create request: %s", err)
		return
	}

	for _, opt := range options {
		opt(req)
	}

	resp, err = c.client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to request: %s", err)
		return
	}

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("err reading body: %s", err)
		return
	}
	return
}
