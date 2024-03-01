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
	BreakLoop:
		for {
			line, prefix, err := body.ReadLine()
			if err != nil {
				if err == io.EOF {
					err = nil
				}
				break BreakLoop
			}
			fmt.Println(prefix, string(line))
		}
	}
}
