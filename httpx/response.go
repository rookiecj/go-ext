package httpx

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	io.Closer
	res         *http.Response
	bufBody     *bufio.Reader
	bodyParsers map[string]BodyParser
}

func (c *Response) Header() map[string][]string {
	return c.res.Header
}

func (c *Response) Status() string {
	return c.res.Status
}

func (c *Response) StatusCode() int {
	return c.res.StatusCode
}

// BufferedReader returns the reader for body
// [Close] will close the body
func (c *Response) BufferedReader() *bufio.Reader {
	return c.bufBody
}

// Closer interface
func (c *Response) Close() {
	// res.Body can be nil with HEAD method
	c.bufBody = nil
	if c.res != nil && c.res.Body != nil {
		c.res.Body.Close()
		c.res.Body = nil
	}
}

func (c *Response) getBodyParser(contentType string) BodyParser {
	lowerType := strings.ToLower(contentType)
	for key, parser := range c.bodyParsers {
		if strings.HasPrefix(lowerType, key) {
			return parser
		}
	}
	return nil
}

// Unmarshal unmarshal body and close the body stream
func (c *Response) Unmarshal(ptrType any) (err error) {

	var contentTypes []string
	var ok bool
	if contentTypes, ok = c.Header()["Content-Type"]; !ok {
		return fmt.Errorf("no content-type found")
	}

	var bodyParser BodyParser
	if bodyParser = c.getBodyParser(contentTypes[0]); bodyParser == nil {
		return fmt.Errorf("no parser found for %s", contentTypes[0])
	}
	// close body here
	defer c.Close()

	// parse body
	err = bodyParser(c.BufferedReader(), ptrType)

	return
}
