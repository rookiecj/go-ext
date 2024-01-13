package httpx

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type Response struct {
	res         *http.Response
	bufBody     *bufio.Reader
	bodyParsers map[string]BodyParser
}

func (resp *Response) Header() map[string][]string {
	return resp.res.Header
}

func (resp *Response) Status() string {
	return resp.res.Status
}

func (resp *Response) StatusCode() int {
	return resp.res.StatusCode
}

// Close will be called on this Response
func (resp *Response) BufferedReader() *bufio.Reader {
	return resp.bufBody
}

// Closer interface
func (resp *Response) Close() {
	// res.Body can be nil with HEAD method
	if resp.res != nil && resp.res.Body != nil {
		resp.res.Body.Close()
	}
}

// unmarshal body and close the body stream
func (resp *Response) Unmarshal(ptrType any) (err error) {

	// parse body
	var data []byte
	if contentTypes, ok := resp.Header()["Content-Type"]; ok {
		if bodyParser := getBodyParser(contentTypes[0]); bodyParser != nil {
			defer resp.Close()
			// read all head
			data, err = io.ReadAll(resp.BufferedReader())
			if err != nil {
				return fmt.Errorf("error while reading: %s", err)
			}

			err = bodyParser(data, ptrType)
			if err == nil {
				return // done
			}
			// fall through
		}
	}

	// fallback to string
	bodyType := reflect.TypeOf(ptrType).Elem()
	if bodyType.Kind() == reflect.String {
		valType := reflect.ValueOf(ptrType)
		valRef := reflect.Indirect(valType)
		valRef.SetString(string(data))
		return nil
	}

	// error
	return fmt.Errorf("unknown error: %s", err)
}
