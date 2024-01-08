package httpx

import (
	"net"
	"net/http"
	"reflect"
	"time"
)

type Builder interface {
	Timeout(timeout time.Duration) Builder
	AddBodyParser(bodyType reflect.Type, parser BodyParser) Builder
	Build() Client
}

type clientBuilder struct {
	timeout     time.Duration
	bodyParsers map[reflect.Type]BodyParser
	// the Transport requests gzip on its own and gets a gzipped response
	disableCompression bool // false
}

func NewBuilder() Builder {
	return &clientBuilder{
		timeout:            DefaultTimeout,
		bodyParsers:        make(map[reflect.Type]BodyParser),
		disableCompression: false,
	}
}

func (c *clientBuilder) Timeout(timeout time.Duration) Builder {
	c.timeout = timeout
	return c
}

func (c *clientBuilder) AddBodyParser(bodyType reflect.Type, parser BodyParser) Builder {
	c.bodyParsers[bodyType] = parser
	return c
}

func (c *clientBuilder) Build() Client {
	// https://go.dev/src/net/http/transport.go
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout: c.timeout,
			//KeepAlive: 15 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: c.timeout,
		DisableCompression:  c.disableCompression,
	}

	client := baseClient{
		client: http.Client{
			Transport: &transport,
			// connection timeout
			Timeout: c.timeout,
		},
		bodyParsers:        c.bodyParsers,
		disableCompression: c.disableCompression,
	}
	return &client
}
