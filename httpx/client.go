package httpx

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultTimeout = 10 * time.Second
)

type Client struct {
	//Do(method string, url url.URL, options ...Option) (res *http.Response, data []byte, err error)
	client http.Client
	//bodyParsers        map[reflect.Type]BodyParser
	bodyParsers        map[string]BodyParser
	disableCompression bool
}

var DefaultClient = NewBuilder().Build()

func (client *Client) Do(method string, url *url.URL, options ...Option) (res *Response, err error) {

	// new request
	req := newRequest()
	for _, opt := range options {
		opt(req)
	}

	// path
	if len(req.path) != 0 {
		url = url.JoinPath(req.path)
	}

	// query
	queries := url.Query()
	for q, values := range req.queries {
		for _, v := range values {
			queries.Add(q, v)
		}
	}
	url.RawQuery = queries.Encode()

	// body
	var body io.Reader = req.body
	if req.body == nil {
		body = nil
	}

	hreq, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		err = fmt.Errorf("failed to create request: %s", err)
		return
	}

	for k, v := range req.headers {
		hreq.Header.Add(k, v)
	}
	if len(req.contentType) != 0 {
		hreq.Header.Set("Content-Type", req.contentType)
	}

	// create http.Request from Request
	resp, err := client.client.Do(hreq)
	if err != nil {
		err = fmt.Errorf("failed to request: %s", err)
		return
	}

	// resp.Body can be nil with HEAD method
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	var data []byte
	if resp.Body != nil {
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			err = fmt.Errorf("err reading body: %s", err)
			return
		}
	}

	res = &Response{
		res:         resp,
		bodyParsers: client.bodyParsers,
		Header:      resp.Header,
		Status:      resp.Status,
		StatusCode:  resp.StatusCode,
		Data:        data,
	}
	return
}

func (client *Client) Get(url *url.URL, options ...Option) (res *Response, err error) {
	res, err = client.Do("GET", url, options...)
	return
}

func (client *Client) Post(url *url.URL, options ...Option) (res *Response, err error) {
	res, err = client.Do("POST", url, options...)
	return
}

func (client *Client) Put(url *url.URL, options ...Option) (res *Response, err error) {
	res, err = client.Do("PUT", url, options...)
	return
}

func (client *Client) Delete(url *url.URL, options ...Option) (res *Response, err error) {
	res, err = client.Do("DELETE", url, options...)
	return
}

func (client *Client) Head(url *url.URL, options ...Option) (res *Response, err error) {
	res, err = client.Do("HEAD", url, options...)
	return
}

//
//func (client *Client)Options(url *url.URL, options ...Option) (res *Response, err error)  {
//	res, err = client.Do("OPTIONS", url, options...)
//	return
//}
//
//func (client *Client)Patch(url *url.URL, options ...Option) (res *Response, err error)  {
//	res, err = client.Do("PATCH", url, options...)
//	return
//}
