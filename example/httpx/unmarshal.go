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
