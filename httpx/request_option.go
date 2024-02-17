package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type ReqOption func(req *Request) error

func newRequest() *Request {
	return &Request{
		headers: make(map[string]string),
		body:    nil,
	}
}

func WithHeader(key, value string) ReqOption {
	return func(req *Request) error {
		if req.headers == nil {
			req.headers = make(map[string]string)
		}
		req.headers[key] = value
		return nil
	}
}

func WithHeaders(headers map[string]string) ReqOption {
	return func(req *Request) error {
		if req.headers == nil {
			req.headers = make(map[string]string)
		}
		for k, v := range headers {
			req.headers[k] = v
		}
		return nil
	}
}

func WithPath(path string) ReqOption {
	return func(req *Request) error {
		req.path += path
		return nil
	}
}

func WithQuery(key, value string) ReqOption {
	return func(req *Request) error {
		if req.queries == nil {
			req.queries = make(map[string][]string)
		}
		req.queries[key] = append(req.queries[key], value)
		return nil
	}
}

func WithQueries(queries map[string][]string) ReqOption {
	return func(req *Request) error {
		if req.queries == nil {
			req.queries = make(map[string][]string)
		}
		for key, values := range queries {
			req.queries[key] = append(req.queries[key], values...)
		}
		return nil
	}
}

func WithBuffer(contentType string, body *bytes.Buffer) ReqOption {
	return func(req *Request) error {
		req.contentType = contentType
		req.body = body
		return nil
	}
}

func WithString(contentType string, body string) ReqOption {
	return WithBytes(contentType, []byte(body))
}

func WithBytes(contentType string, body []byte) ReqOption {
	return WithBuffer(contentType, bytes.NewBuffer(body))
}

func WithJsonString(json string) ReqOption {
	return WithString("application/json; charset=UTF-8", json)
}

func WithJsonObject(body any) ReqOption {
	return func(req *Request) error {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		req.contentType = "application/json; charset=UTF-8"
		req.body = bytes.NewBuffer(b)
		return nil
	}
}

func WithFormData(fields map[string]string) ReqOption {
	return func(req *Request) error {
		if req.body == nil {
			req.body = &bytes.Buffer{}
		}

		req.contentType = "application/x-www-form-urlencoded"
		data := url.Values{}
		for key, value := range fields {
			data.Add(key, value)
		}
		req.body.WriteString(data.Encode())
		return nil
	}
}

func WithMultipartFile(fieldName string, file *os.File) ReqOption {
	return WithMultipartReader(fieldName, file.Name(), file)
}

func WithMultipartReader(fieldName string, filename string, reader io.Reader) ReqOption {
	return func(req *Request) error {
		if req.body == nil {
			req.body = &bytes.Buffer{}
		}

		mw := multipart.NewWriter(req.body)
		defer func() {
			if mw != nil {
				mw.Close()
			}
		}()

		formWriter, err := mw.CreateFormFile(fieldName, filename)
		if err != nil {
			return fmt.Errorf("err create form file: %s", err)
		}
		if _, err = io.Copy(formWriter, reader); err != nil {
			return fmt.Errorf("copy error: %s", err)
		}
		req.contentType = mw.FormDataContentType()
		return nil
	}
}
