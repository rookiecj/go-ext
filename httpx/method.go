package httpx

import (
	"fmt"
	"net/url"
	"reflect"
)

func fetch[Body any](client *Client, method string, url *url.URL, options ...Option) (body Body, err error) {
	if client == nil {
		client = DefaultClient
	}

	resp, err := client.Do(method, url, options...)
	if err != nil {
		err = fmt.Errorf("failed to request: %s", err)
		return
	}

	// parse body
	if contentTypes, ok := resp.Header["Content-Type"]; ok {
		if bodyParser := getBodyParser(contentTypes[0]); bodyParser != nil {
			err = bodyParser(resp.Data, &body)
			if err == nil {
				return
			}
		}
	}

	// string
	bodyType := reflect.TypeOf(body)
	if bodyType.Kind() == reflect.String {
		valType := reflect.ValueOf(body)
		valRef := reflect.Indirect(valType)
		valRef.SetString(string(resp.Data))
		err = nil
		return
	}

	err = fmt.Errorf("parse error: %s", err)
	return
}

func Get[Body any](client *Client, url *url.URL, options ...Option) (body Body, err error) {
	body, err = fetch[Body](client, "GET", url, options...)
	return
}

func Post[Body any](client *Client, url *url.URL, options ...Option) (body Body, err error) {
	body, err = fetch[Body](client, "POST", url, options...)
	return
}

func Put[Body any](client *Client, url *url.URL, options ...Option) (body Body, err error) {
	body, err = fetch[Body](client, "PUT", url, options...)
	return
}

func Delete[Body any](client *Client, url *url.URL, options ...Option) (body Body, err error) {
	body, err = fetch[Body](client, "PUT", url, options...)
	return
}

func Head(client *Client, url *url.URL, options ...Option) (body *any, err error) {
	_, err = fetch[any](client, "HEAD", url, options...)
	return
}

//
//func Options[Body any](client *Client, url *url.URL, options ...Option) (body *Body, err error) {
//	body, err = fetch[Body](client, "OPTIONS", url, options...)
//	return
//}
//
//func Patch[Body any](client *Client, url *url.URL, options ...Option) (body *Body, err error) {
//	body, err = fetch[Body](client, "PATCH", url, options...)
//	return
//}
