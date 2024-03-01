package main

import (
	"bytes"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
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

	req, err := proto.Marshal(request)
	if err != nil {
		log.Fatalf("Unable to marshal request : %v", err)
	}

	resp, err := http.Post("http://0.0.0.0:8080/post", "application/x-binary", bytes.NewReader(req))
	if err != nil {
		log.Fatalf("Unable to read from the server : %v", err)
	}
	defer resp.Body.Close()

	respObj := &dto.PostResponse{}
	pb.ProtoBufBodyParser(resp.Body, respObj)
	return respObj
}

func main() {

	request := &dto.PostRequest{Id: 1}
	resp := makeRequest(request)
	fmt.Printf("Response from API is : %v\n", resp.Body)
}
