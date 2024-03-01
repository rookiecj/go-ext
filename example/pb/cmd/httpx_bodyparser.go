package main

import (
	"fmt"
	"github.com/rookiecj/go-langext/httpx"
	"google.golang.org/protobuf/proto"
	"log"
	pb "protobuf-example"
	"protobuf-example/proto/dto"
)

type testPost struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

var testPostUrl = "http://0.0.0.0:8080/post"

func makeRequest(request *dto.PostRequest) *dto.PostResponse {

	c := httpx.NewClient(httpx.WithBodyParser("application/protobuf", pb.ProtoBufBodyParser))

	req, err := proto.Marshal(request)
	if err != nil {
		log.Fatalf("Unable to marshal request : %v", err)
	}
	resp, err := c.Post(testPostUrl,
		httpx.WithBytes("application/x-binary", req))
	if err != nil {
		log.Fatalf("Unable to read from the server : %v", err)
	}
	defer resp.Close()

	respObj := &dto.PostResponse{}
	err = resp.Unmarshal(respObj)
	return respObj
}

func main() {

	request := &dto.PostRequest{Id: 1}
	resp := makeRequest(request)
	fmt.Printf("Response is : %v\n", resp)
}
