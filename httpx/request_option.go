package httpx

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"
	"os"
)

type Request struct {
	// fields can be hidden
	// becuase the request only be created and accessed in the client
	headers map[string]string
	path    string
	queries map[string][]string

	contentType string
	body        *bytes.Buffer
}

type Option func(req *Request)

func newRequest() *Request {
	return &Request{
		headers: make(map[string]string),
		body:    nil,
	}
}

func WithHeader(key, value string) func(req *Request) {
	return func(req *Request) {
		if req.headers == nil {
			req.headers = make(map[string]string)
		}
		req.headers[key] = value
	}
}

func WithPath(path string) func(req *Request) {
	return func(req *Request) {
		req.path += path
	}
}

func WithQuery(key, value string) func(req *Request) {
	return func(req *Request) {
		if req.queries == nil {
			req.queries = make(map[string][]string)
		}
		req.queries[key] = append(req.queries[key], value)
	}
}

func WithQueries(queries map[string][]string) func(req *Request) {
	return func(req *Request) {
		if req.queries == nil {
			req.queries = make(map[string][]string)
		}
		for key, values := range queries {
			req.queries[key] = append(req.queries[key], values...)
		}
	}
}

func WithBuffer(contentType string, body *bytes.Buffer) func(req *Request) {
	return func(req *Request) {
		req.contentType = contentType
		req.body = body
	}
}

func WithString(contentType string, body string) func(req *Request) {
	return WithBytes(contentType, []byte(body))
}

func WithBytes(contentType string, body []byte) func(req *Request) {
	return WithBuffer(contentType, bytes.NewBuffer(body))
}

func WithJsonString(json string) func(req *Request) {
	return WithString("application/json; charset=UTF-8", json)
}

func WithJsonObject(body any) func(req *Request) {
	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return WithJsonString(string(b))
}

func WithFormData(fields map[string]string) func(req *Request) {
	return func(req *Request) {
		if req.body == nil {
			req.body = &bytes.Buffer{}
		}

		req.contentType = "application/x-www-form-urlencoded"
		data := url.Values{}
		for key, value := range fields {
			data.Add(key, value)
		}
		req.body.WriteString(data.Encode())
	}
}

func WithFile(fieldName string, file *os.File) func(req *Request) {
	return func(req *Request) {
		if req.body == nil {
			req.body = &bytes.Buffer{}
		}

		mw := multipart.NewWriter(req.body)
		defer func() {
			if mw != nil {
				mw.Close()
			}
		}()

		formWriter, err := mw.CreateFormFile(fieldName, file.Name())
		if err != nil {
			panic(err)
		}
		if _, err = io.Copy(formWriter, file); err != nil {
			panic(err)
		}
		req.contentType = mw.FormDataContentType()
	}
}
