package main

import (
	"fmt"
	"github.com/rookiecj/go-langext/httpx"
)

var testPostUrl = "https://jsonplaceholder.typicode.com"

type testPost struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

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
