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

	newPost := testPost{
		UserId: 1,
		Title:  "New title",
		Body:   "New Body",
	}

	if res, err := client.Post(testPostUrl,
		httpx.WithPath("/posts"),
		httpx.WithJsonObject(newPost),
	); err == nil {
		defer res.Close()

		var resPost testPost
		if err := res.Unmarshal(&resPost); err == nil {
			fmt.Printf("res: %v\n", resPost)
		}
	}
}

```

## Extensions

There are two extentions:

- lang extensions
- http extensions

### lang extensions

- [X] container: priority queue
- [X] container: sorted map
- [X] container: set, sorted set

### http extensions

- [X] `NewClient` http client with option
- [X] `BodyParser` Http response body parser

## Todo

- [ ] http response adapter for `go-stream`
