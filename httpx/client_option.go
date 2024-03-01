package httpx

import (
	"time"
)

const (
	DefaultTimeout = 10 * time.Second
)

type clientOptions struct {
	timeout time.Duration
	//bodyParsers map[reflect.Type]BodyParser
	bodyParsers map[string]BodyParser
	// the Transport requests gzip on its own and gets a gzipped response
	disableCompression bool // false
	headers            map[string][]string
}

type ClientOption func(clientOptions *clientOptions)

func defaultClientOptions() clientOptions {
	co := clientOptions{
		timeout: DefaultTimeout,
		//bodyParsers:        make(map[reflect.Type]BodyParser),
		bodyParsers:        make(map[string]BodyParser),
		disableCompression: false,
	}
	co.bodyParsers["application/json"] = JsonBodyParser
	return co
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(clientOptions *clientOptions) {
		clientOptions.timeout = timeout
	}
}

func WithDefaultHeaders(headers map[string][]string) ClientOption {
	return func(clientOptions *clientOptions) {
		clientOptions.headers = headers
	}
}

func WithBodyParser(contentType string, parser BodyParser) ClientOption {
	return func(clientOptions *clientOptions) {
		clientOptions.bodyParsers[contentType] = parser
	}
}
