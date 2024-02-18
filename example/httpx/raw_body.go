package main

import (
	"fmt"
	"github.com/rookiecj/go-langext/httpx"
	"io"
)

var testPostUrl = "https://jsonplaceholder.typicode.com"

func main() {

	client := httpx.NewClient()

	if res, err := client.Get(testPostUrl, httpx.WithPath("/posts/1")); err == nil {
		defer res.Close()
		body := res.BufferedReader()
		for line, prefix, err := body.ReadLine(); err == nil || err == io.EOF; line, prefix, err = body.ReadLine() {
			fmt.Println(prefix, string(line))
			if err == io.EOF {
				break
			}
		}
	}
}
