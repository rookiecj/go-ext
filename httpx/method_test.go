package httpx

import (
	"fmt"
	"net/url"
	"testing"
)

type testTodo struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Done   bool
}

const (
	testTodoUrl = "https://jsonplaceholder.typicode.com/posts"
)

var testClient = NewBuilder().Build()

func TestFetchSync(t *testing.T) {
	type args[Body any] struct {
		method string
		url    url.URL
	}
	type testCase[Body any] struct {
		name    string
		client  Client
		args    args[Body]
		wantErr error
	}
	tests := []testCase[testTodo]{
		{
			name:   "get todo 1",
			client: testClient,
			args: args[testTodo]{
				method: "GET",
				url: func() url.URL {
					u, _ := url.Parse(testTodoUrl + "/1")
					return *u
				}(),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if body, err := fetch[testTodo](tt.client, tt.args.method, tt.args.url); err != nil {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				fmt.Println(body)
			}
		})
	}
}
