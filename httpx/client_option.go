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
