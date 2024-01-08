package httpx

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
)

type GetOption func(r *GetRequest)
type HeadOption func(r *HeadRequest)
type PostOption func(r *PostRequest)
type PutOption func(r *PutRequest)
type DeleteOption func(r *DeleteRequest)
type OptionsOption func(r *OptionsRequest)
type PatchOption func(r *PatchRequest)

type GetRequest struct {
	Header map[string][]string
}

type HeadRequest struct {
	Header map[string][]string
}

type PostRequest struct {
	Header map[string][]string
}

type PutRequest struct {
	Header map[string][]string
}

type DeleteRequest struct {
	Header map[string][]string
}

type OptionsRequest struct {
	Header map[string][]string
}

type PatchRequest struct {
	Header map[string][]string
}

func fetch[Body any](client Client, method string, url url.URL, options ...RequestOption) (body Body, err error) {

	// new client and Do
	resp, data, err := client.Do(method, url, options...)
	if err != nil {
		err = fmt.Errorf("failed to request: %s", err)
		return
	}

	// parse body
	contentTypes := resp.Header["Content-Type"]
	if bodyParser := getBodyParser(contentTypes[0]); bodyParser != nil {
		err = bodyParser(data, &body)
		if err == nil {
			return
		}
	}

	// string
	bodyType := reflect.TypeOf(body)
	if bodyType.Kind() == reflect.String {
		valType := reflect.ValueOf(body)
		valRef := reflect.Indirect(valType)
		valRef.SetString(string(data))
		err = nil
		return
	}

	err = fmt.Errorf("parse error: %s", err)
	return
}

func Get[Body any](client Client, url url.URL, options ...GetOption) (body Body, err error) {
	var req GetRequest

	for _, opt := range options {
		opt(&req)
	}

	// GetOption -> RequestOption
	reqOption := func(httpReq *http.Request) {
		httpReq.Header = req.Header
	}
	return fetch[Body](client, "GET", url, reqOption)
}

func Head[Body any](client Client, url url.URL, options ...HeadOption) (body Body, err error) {
	var req HeadRequest
	for _, opt := range options {
		opt(&req)
	}

	reqOption := func(httpReq *http.Request) {
		httpReq.Header = req.Header
	}
	return fetch[Body](client, "HEAD", url, reqOption)
}

func Post[Body any](client Client, url url.URL, options ...PostOption) (body Body, err error) {
	var req PostRequest
	for _, opt := range options {
		opt(&req)
	}

	reqOption := func(httpReq *http.Request) {
		httpReq.Header = req.Header
	}
	return fetch[Body](client, "POST", url, reqOption)
}

func Put[Body any](client Client, url url.URL, options ...PutOption) (body Body, err error) {
	var req PutRequest
	for _, opt := range options {
		opt(&req)
	}

	reqOption := func(httpReq *http.Request) {
		httpReq.Header = req.Header
	}
	return fetch[Body](client, "PUT", url, reqOption)
}

func Delete[Body any](client Client, url url.URL, options ...DeleteOption) (body Body, err error) {
	var req DeleteRequest
	for _, opt := range options {
		opt(&req)
	}

	reqOption := func(httpReq *http.Request) {
		httpReq.Header = req.Header
	}
	return fetch[Body](client, "DELETE", url, reqOption)
}

func Options[Body any](client Client, url url.URL, options ...OptionsOption) (body Body, err error) {
	var req OptionsRequest
	for _, opt := range options {
		opt(&req)
	}

	reqOption := func(httpReq *http.Request) {
		httpReq.Header = req.Header
	}
	return fetch[Body](client, "OPTIONS", url, reqOption)
}

func Patch[Body any](client Client, url url.URL, options ...PatchOption) (body Body, err error) {
	var req PatchRequest
	for _, opt := range options {
		opt(&req)
	}

	reqOption := func(httpReq *http.Request) {
		httpReq.Header = req.Header
	}
	return fetch[Body](client, "PATCH", url, reqOption)
}
