# Go Extension functions

`go-langext` is a collection of complementary functions which is missing in Go library.

## How to install

```
go get github.com/rookiecj/go-langext
```

## How to use

Here is httpx example:

```go
package main

import (
	"fmt"
	"github.com/rookiecj/go-langext/httpx"
)

type testPost struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

var testPostUrl = "https://jsonplaceholder.typicode.com"

func main() {
	client := httpx.NewClient()

	var res *httpx.Response
	var err error
	if res, err = client.Get(testPostUrl, httpx.WithPath("/posts/1")); err != nil {
		panic(err)
	}

	post := testPost{}
	if err = res.Unmarshal(&post); err != nil {
		panic(err)
	}
	fmt.Printf("post: %v\n", post)
	res.Close()
}


```

## Extensions

There are two extentions:

- lang extensions
- http extensions

### lang extensions

- [X] container: priority queue
- [X] sorted map

### http extensions

- [X] http client
- [X] Http response body parser

## Todo

- [ ] http response adapter for `go-stream`
