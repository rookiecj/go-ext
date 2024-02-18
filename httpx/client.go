package httpx

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	neturl "net/url"
)

type Client struct {
	client             http.Client
	defaultHeaders     map[string][]string
	bodyParsers        map[string]BodyParser
	disableCompression bool
}

var DefaultClient = NewClient()

func NewClient(options ...ClientOption) *Client {

	co := defaultClientOptions()
	for _, opt := range options {
		opt(&co)
	}

	// https://go.dev/src/net/http/transport.go
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout: co.timeout,
			//KeepAlive: 15 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: co.timeout,
		DisableCompression:  co.disableCompression,
	}

	client := Client{
		client: http.Client{
			Transport: &transport,
			// connection timeout
			Timeout: co.timeout,
		},
		bodyParsers:        co.bodyParsers,
		disableCompression: co.disableCompression,
		defaultHeaders:     co.headers,
	}
	return &client
}

func (c *Client) Do(method string, url *url.URL, options ...ReqOption) (res *Response, err error) {

	// new request
	req := newRequest()

	// options
	for _, option := range options {
		if err = option(req); err != nil {
			return
		}
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

	// headers
	// add default headers
	for k, values := range c.defaultHeaders {
		for _, v := range values {
			hreq.Header.Add(k, v)
		}
	}
	// request headers
	for k, values := range req.headers {
		for _, v := range values {
			hreq.Header.Add(k, v)
		}
	}
	if len(req.contentType) != 0 {
		hreq.Header.Set("Content-Type", req.contentType)
	}

	// make request
	resp, err := c.client.Do(hreq)
	if err != nil {
		err = fmt.Errorf("failed to request: %s", err)
		return
	}

	// response
	res = &Response{
		res:         resp,
		bufBody:     nil,
		bodyParsers: c.bodyParsers,
	}

	// 204 No Content
	if resp.StatusCode == 204 {
		if resp.Body != nil {
			resp.Body.Close()
			resp.Body = nil
		}
	} else {
		res.bufBody = bufio.NewReaderSize(resp.Body, 4*1024)
	}
	return
}

func (c *Client) Get(url string, options ...ReqOption) (res *Response, err error) {
	newUrl, err := neturl.Parse(url)
	if err != nil {
		err = fmt.Errorf("parse error: %v", err)
		return
	}
	res, err = c.Do("GET", newUrl, options...)
	return
}

func (c *Client) Post(url string, options ...ReqOption) (res *Response, err error) {
	newUrl, err := neturl.Parse(url)
	if err != nil {
		err = fmt.Errorf("parse error: %v", err)
		return
	}

	res, err = c.Do("POST", newUrl, options...)
	return
}

func (c *Client) Put(url string, options ...ReqOption) (res *Response, err error) {
	newUrl, err := neturl.Parse(url)
	if err != nil {
		err = fmt.Errorf("parse error: %v", err)
		return
	}
	res, err = c.Do("PUT", newUrl, options...)
	return
}

//func (client *Client) Patch(url string, headers map[string]string, jsonObj any) (res *Response, err error) {
//	newUrl, err := neturl.Parse(url)
//	if err != nil {
//		err = fmt.Errorf("parse error: %v", err)
//		return
//	}
//	res, err = client.PatchWith(newUrl, WithHeaders(headers), WithJsonObject(jsonObj))
//	return
//}
//
//func (client *Client) PatchWith(url *url.URL, options ...ReqOption) (res *Response, err error) {
//	res, err = client.Do("PATCH", url, options...)
//	return
//}

func (c *Client) Delete(url string, options ...ReqOption) (res *Response, err error) {
	newUrl, err := neturl.Parse(url)
	if err != nil {
		err = fmt.Errorf("parse error: %v", err)
		return
	}
	res, err = c.Do("DELETE", newUrl, options...)
	return
}

func (c *Client) Head(url string, options ...ReqOption) (res *Response, err error) {
	newUrl, err := neturl.Parse(url)
	if err != nil {
		err = fmt.Errorf("parse error: %v", err)
		return
	}
	res, err = c.Do("HEAD", newUrl, options...)
	// header has no response body, close here
	res.Close()
	return
}

func (c *Client) Options(url string, options ...ReqOption) (res *Response, err error) {
	newUrl, err := neturl.Parse(url)
	if err != nil {
		err = fmt.Errorf("parse error: %v", err)
		return
	}
	res, err = c.Do("OPTIONS", newUrl, options...)
	// Options has no response body, close here
	res.Close()
	return
}
