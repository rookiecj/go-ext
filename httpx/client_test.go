package httpx

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rookiecj/go-langext/langx"
	"github.com/rookiecj/go-langext/mapper"
	"io"
	"mime/multipart"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

const (
	testPostUrl = "https://jsonplaceholder.typicode.com"
)

func TestGetSimple(t *testing.T) {
	type args struct {
		client  *Client
		url     string
		options []ReqOption
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
				client: NewClient(),
				url: func() string {
					//u, _ := url.Parse(testPostUrl)
					//return u
					return testPostUrl
				}(),
				options: []ReqOption{
					WithPath("/posts/1"),
				},
			},
			wantBody: testPosts[0],
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.args.client.Get(tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var gotBody testPost
			gotErr := gotResp.Unmarshal(&gotBody)
			if gotErr != nil {
				t.Errorf("Get() error while unmarshal resp %v", gotErr)
			} else if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("Get() gotBody = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestGetMulti(t *testing.T) {
	type args struct {
		client  *Client
		url     string
		options []ReqOption
	}
	type testCase[Body any] struct {
		name     string
		args     args
		wantBody []Body
		wantErr  bool
	}
	tests := []testCase[testPost]{
		{
			name: "GET /posts",
			args: args{
				client: NewClient(),
				url: func() string {
					//u, _ := url.Parse(testPostUrl)
					//return u
					return testPostUrl
				}(),
				options: []ReqOption{
					WithPath("/posts"),
				},
			},
			wantBody: testPosts,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.args.client.Get(tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMulti() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotBody := []testPost{}
			gotErr := gotResp.Unmarshal(&gotBody)
			if gotErr != nil {
				t.Errorf("GetMulti() error while unmarshal resp %v", gotErr)
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
		url     string
		options []ReqOption
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
				client: NewClient(),
				url: func() string {
					//u, _ := url.Parse(testPostUrl)
					//return u
					return testPostUrl
				}(),
				options: []ReqOption{
					WithPath("/posts/1/comments"),
				},
			},
			wantBody: testComments,
			wantErr:  false,
		},
		{
			name: "GET - /comments?postId=1",
			args: args{
				client: NewClient(),
				url: func() string {
					//u, _ := url.Parse(testPostUrl)
					//return u
					return testPostUrl
				}(),
				options: []ReqOption{
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
			gotResp, err := tt.args.client.Get(tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMultiComments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotBody := []testComment{}
			gotErr := gotResp.Unmarshal(&gotBody)
			if gotErr != nil {
				t.Errorf("GetMultiComments() error while unmarshal resp %v", gotErr)
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
		url     string
		options []ReqOption
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
				client: NewClient(),
				url: func() string {
					//u, _ := url.Parse(testPostUrl)
					//return u
					return testPostUrl
				}(),
				options: []ReqOption{
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
			gotResp, err := tt.args.client.Post(tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotBody := testPostResponse{}
			gotErr := gotResp.Unmarshal(&gotBody)
			if gotErr != nil {
				t.Errorf("Post() error while unmarshal resp %v", gotErr)
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
		url     string
		options []ReqOption
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
				client: NewClient(),
				url: func() string {
					return testPostUrl
				}(),
				options: []ReqOption{
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
			gotResp, err := tt.args.client.Post(tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotBody := testPostFormResponse{}
			gotErr := gotResp.Unmarshal(&gotBody)
			if gotErr != nil {
				t.Errorf("PostForm() error while unmarshal resp %v", gotErr)
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("PostForm() gotBody = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestPostFileSimple(t *testing.T) {
	type args struct {
		client  *Client
		url     string
		options []ReqOption
	}
	type testCase struct {
		name     string
		args     args
		wantBody any
		wantErr  bool
	}

	serverUrl := startServer()

	tests := []testCase{
		{
			name: "POST - WithMultipartFile",
			args: args{
				client: NewClient(),
				url: func() string {
					u, _ := url.Parse(serverUrl)
					return u.String()
				}(),
				options: []ReqOption{
					WithPath("/posts"),
					WithMultipartReader("testfieldname", "testfilename.txt", func() io.Reader {
						msg := "this is file"
						b := bytes.NewBufferString(msg)
						return b
					}()),
				},
			},
			wantBody: "this is file",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.args.client.Post(tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post File() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			//--b7fd4733b5b348e64082bc8bee0f38de8ea6c5a28c985cff6c83cff24500
			//Content-Disposition: form-data; name="testfieldname"; filename="testfilename"
			//Content-Type: application/octet-stream
			//
			//this is file
			//--b7fd4733b5b348e64082bc8bee0f38de8ea6c5a28c985cff6c83cff24500--
			//

			// first line은 boundary
			bufBody := gotResp.BufferedReader()
			size := 2
			boundary := ""
			thunk, err := bufBody.Peek(size)
			for err == nil {
				if strings.HasSuffix(string(thunk), "\r\n") {
					boundary = string(thunk)
					break
				}
				size += 1
				thunk, err = bufBody.Peek(size)
			}
			boundary = strings.TrimSpace(boundary)
			boundary = strings.Trim(boundary, "--")

			partReader := multipart.NewReader(bufBody, boundary)
			part, err := partReader.NextPart()
			gotBody := ""
			for err == nil {
				for k, v := range part.Header {
					fmt.Printf("part.Header: %s = '%v'\n", k, v)
				}
				fmt.Printf("part.FileName: %v\n", part.FileName())
				var buff = make([]byte, 0, 64)
				n, nerr := part.Read(buff[len(buff):cap(buff)])
				if nerr != nil && !errors.Is(nerr, io.EOF) {
					panic(nerr)
				}
				buff = buff[:len(buff)+n]
				gotBody = string(buff)
				fmt.Printf("part.Read: %d, %s\n", n, gotBody)
				part.Close()
				part, err = partReader.NextPart()
			}
			gotResp.Close()

			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("Post File() gotBody = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
	stopServer()
}

func TestPutSimple(t *testing.T) {
	type args struct {
		client  *Client
		url     string
		options []ReqOption
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
				client: NewClient(),
				url: func() string {
					u, _ := url.Parse(testPostUrl)
					return u.String()
				}(),
				options: []ReqOption{
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
			gotResp, err := tt.args.client.Put(tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotBody := testPutResponse{}
			gotErr := gotResp.Unmarshal(&gotBody)
			if gotErr != nil {
				t.Errorf("Put() error while unmarshal resp %v", gotErr)
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
		url     string
		options []ReqOption
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
				client: NewClient(),
				url: func() string {
					u, _ := url.Parse(testPostUrl)
					return u.String()
				}(),
				options: []ReqOption{
					WithPath("/posts/1"),
				},
			},
			wantBody: testDeleteResponse{
				//Id: 1, // empty
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.args.client.Delete(tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotBody := testDeleteResponse{}
			gotErr := gotResp.Unmarshal(&gotBody)
			if gotErr != nil {
				t.Errorf("Delete() error while unmarshal resp %v", gotErr)
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("Delete() gotBody = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestClient_Head(t *testing.T) {

	type args struct {
		url     string
		options []ReqOption
	}
	tests := []struct {
		name        string
		client      *Client
		args        args
		wantHeaders []string
		wantErr     bool
	}{
		{
			name:   "Head",
			client: NewClient(),
			args: args{
				url: func() string {
					u, _ := url.Parse(testPostUrl)
					return u.String()
				}(),
				options: []ReqOption{},
			},
			wantHeaders: []string{
				"Date",
				"Connection",
				"Cache-Control",
				"Age",
				"Server",
				"Content-Type",
				"Via",
				"Last-Modified",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.client.Head(tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Head() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotValues := gotRes.Header()
			for _, header := range tt.wantHeaders {
				if _, ok := gotValues[header]; !ok {
					t.Errorf("Head(): want %v, got %v", header, gotValues)
				}
			}
		})
	}
}

func TestClient_Options(t *testing.T) {

	type args struct {
		url     string
		options []ReqOption
	}
	tests := []struct {
		name        string
		client      *Client
		args        args
		wantHeaders map[string][]string
		wantErr     bool
	}{
		{
			name:   "Options",
			client: NewClient(),
			args: args{
				url: func() string {
					u, _ := url.Parse("https://example.org")
					return u.String()
				}(),
				options: []ReqOption{},
			},
			wantHeaders: func() map[string][]string {
				headers := make(map[string][]string)
				headers["Allow"] = []string{"GET", "POST", "HEAD", "OPTIONS"}
				return headers
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.client.Options(tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Options() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for key, values := range tt.wantHeaders {
				gotValues := gotRes.Header()[key]
				//Header = {http.Header}
				//0 = Allow -> len:1, cap:1
				//key = {string} "Allow"
				//value = {[]string} len:1, cap:1
				//0 = {string} "OPTIONS, GET, HEAD, POST"
				for _, value := range values {
					//if !slices.Contains(gotValues, value) {
					//	t.Errorf("Options(): want %v, got %v", values, gotValues)
					//}
					if !strings.Contains(gotValues[0], value) {
						t.Errorf("Options(): want %v, got %v", values, gotValues)
					}
				}
			}
		})
	}
}
