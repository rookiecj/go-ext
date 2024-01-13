package httpx

import (
	"fmt"
	"github.com/rookiecj/go-langext/langx"
	"github.com/rookiecj/go-langext/mapper"
	"net/url"
	"reflect"
	"testing"
)

const (
	testPostUrl = "https://jsonplaceholder.typicode.com"
)

var testClient = NewBuilder().Build()

func TestFetch(t *testing.T) {
	type args[Body any] struct {
		method string
		url    *url.URL
	}
	type testCase[Body any] struct {
		name    string
		client  *Client
		args    args[Body]
		wantErr error
	}
	tests := []testCase[testTodo]{
		{
			name:   "GET 1",
			client: testClient,
			args: args[testTodo]{
				method: "GET",
				url: func() *url.URL {
					u, _ := url.Parse(testPostUrl + "/posts/1")
					return u
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
				fmt.Printf("%+v\n", body)
			}
		})
	}
}

func TestGetSimple(t *testing.T) {
	type args struct {
		client  *Client
		url     *url.URL
		options []Option
	}
	type testCase[Body any] struct {
		name     string
		args     args
		wantBody Body
		wantErr  bool
	}
	tests := []testCase[testPost]{
		{
			name: "GET - WithPath 1",
			args: args{
				client: NewBuilder().Build(),
				url: func() *url.URL {
					u, _ := url.Parse(testPostUrl)
					return u
				}(),
				options: []Option{
					WithPath("/posts/1"),
				},
			},
			wantBody: testPosts[0],
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := Get[testPost](tt.args.client, tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("Get() gotBody = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestGetMulti(t *testing.T) {
	type args struct {
		client  *Client
		url     *url.URL
		options []Option
	}
	type testCase[Body any] struct {
		name     string
		args     args
		wantBody []Body
		wantErr  bool
	}
	tests := []testCase[testPost]{
		{
			name: "GET",
			args: args{
				client: NewBuilder().Build(),
				url: func() *url.URL {
					u, _ := url.Parse(testPostUrl)
					return u
				}(),
				options: []Option{
					WithPath("/posts"),
				},
			},
			wantBody: testPosts,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := Get[[]testPost](tt.args.client, tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMulti() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotBody) != len(tt.wantBody) {
				t.Errorf("GetMulti() len want %d, got %d", len(tt.wantBody), len(gotBody))
			}

			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("GetMulti() gotBody = %v", gotBody)
				t.Errorf("GetMulti() want %v", tt.wantBody)
			}
		})
	}
}

func TestGetMultiComments(t *testing.T) {
	type args struct {
		client  *Client
		url     *url.URL
		options []Option
	}
	type testCase[Body any] struct {
		name     string
		args     args
		wantBody []Body
		wantErr  bool
	}
	tests := []testCase[testComment]{
		{
			name: "GET - /posts/1/comments",
			args: args{
				client: NewBuilder().Build(),
				url: func() *url.URL {
					u, _ := url.Parse(testPostUrl)
					return u
				}(),
				options: []Option{
					WithPath("/posts/1/comments"),
				},
			},
			wantBody: testComments,
			wantErr:  false,
		},
		{
			name: "GET - /comments?postId=1",
			args: args{
				client: NewBuilder().Build(),
				url: func() *url.URL {
					u, _ := url.Parse(testPostUrl)
					return u
				}(),
				options: []Option{
					WithPath("/comments"),
					WithQuery("postId", "1"),
				},
			},
			wantBody: testComments,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := Get[[]testComment](tt.args.client, tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMultiComments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotBody) != len(tt.wantBody) {
				t.Errorf("GetMultiComments() len want %d, got %d", len(tt.wantBody), len(gotBody))
			}

			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("GetMultiComments() gotBody = %v", gotBody)
				t.Errorf("GetMultiComments() want %v", tt.wantBody)
			}
		})
	}
}

func TestPostSimple(t *testing.T) {
	type args struct {
		client  *Client
		url     *url.URL
		options []Option
	}
	type testCase struct {
		name     string
		args     args
		wantBody any
		wantErr  bool
	}
	tests := []testCase{
		{
			name: "POST - WithJsonObject 1",
			args: args{
				client: NewBuilder().Build(),
				url: func() *url.URL {
					u, _ := url.Parse(testPostUrl)
					return u
				}(),
				options: []Option{
					WithPath("/posts"),
					WithJsonObject(testPosts[0]),
				},
			},
			wantBody: func() testPostResponse {
				var res testPostResponse

				// via Mapper
				mp := mapper.NewMapperWithTag("fieldname")
				//var src = testPosts[0]
				//mp.Map(res, src)
				//res.Id = 101

				// via copy and set
				dupPost := langx.Copy(testPosts[0], func(dup *testPost) {
					dup.Id = 101
				})
				// map
				err := mp.Map(&res, dupPost)
				if err != nil {
					panic(err)
				}
				return res
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := Post[testPostResponse](tt.args.client, tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("Post() gotBody = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestPostFormSimple(t *testing.T) {
	type args struct {
		client  *Client
		url     *url.URL
		options []Option
	}
	type testCase struct {
		name     string
		args     args
		wantBody any
		wantErr  bool
	}
	tests := []testCase{
		{
			name: "POST - WithFormData 1",
			args: args{
				client: NewBuilder().Build(),
				url: func() *url.URL {
					u, _ := url.Parse(testPostUrl)
					return u
				}(),
				options: []Option{
					WithPath("/posts"),
					WithFormData(func() map[string]string {
						//	type testPost struct {
						//	UserId int    `json:"userId"`
						//	Id     int    `json:"id"`
						//	Title  string `json:"title"`
						//	Body   string `json:"body"`
						//}
						var fields = make(map[string]string)
						fields["userId"] = "1"
						fields["id"] = "1"
						fields["title"] = "form title"
						fields["body"] = "form body"
						return fields
					}()),
				},
			},
			wantBody: testPostFormResponse{
				Id:     101,
				UserId: 1,
				Title:  "form title",
				Body:   "form body",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := Post[testPostFormResponse](tt.args.client, tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("PostForm() gotBody = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestPutSimple(t *testing.T) {
	type args struct {
		client  *Client
		url     *url.URL
		options []Option
	}
	type testCase[ResBody any] struct {
		name     string
		args     args
		wantBody ResBody
		wantErr  bool
	}
	tests := []testCase[testPutResponse]{
		{
			name: "PUT - WithJsonObject 1",
			args: args{
				client: NewBuilder().Build(),
				url: func() *url.URL {
					u, _ := url.Parse(testPostUrl)
					return u
				}(),
				options: []Option{
					WithPath("/posts/1"),
					WithJsonObject(testPosts[0]),
				},
			},
			wantBody: testPutResponse{
				Id: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := Put[testPutResponse](tt.args.client, tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("Put() gotBody = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestDeleteSimple(t *testing.T) {
	type args struct {
		client  *Client
		url     *url.URL
		options []Option
	}
	type testCase[ResBody any] struct {
		name     string
		args     args
		wantBody ResBody
		wantErr  bool
	}
	tests := []testCase[testDeleteResponse]{
		{
			name: "DELETE - WithJsonObject 1",
			args: args{
				client: NewBuilder().Build(),
				url: func() *url.URL {
					u, _ := url.Parse(testPostUrl)
					return u
				}(),
				options: []Option{
					WithPath("/posts/1"),
				},
			},
			wantBody: testDeleteResponse{
				Id: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := Delete[testDeleteResponse](tt.args.client, tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("Delete() gotBody = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}
