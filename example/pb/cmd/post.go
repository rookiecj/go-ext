package main

import (
	"bytes"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
	"protobuf-example/proto/dto"
)

func makeRequest(request *dto.PostRequest) *dto.PostResponse {

	req, err := proto.Marshal(request)
	if err != nil {
		log.Fatalf("Unable to marshal request : %v", err)
	}

	resp, err := http.Post("http://0.0.0.0:8080/post", "application/x-binary", bytes.NewReader(req))
	if err != nil {
		log.Fatalf("Unable to read from the server : %v", err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Unable to read bytes from request : %v", err)
	}

	respObj := &dto.PostResponse{}
	proto.Unmarshal(respBytes, respObj)
	return respObj
}

func main() {

	request := &dto.PostRequest{Id: 1}
	resp := makeRequest(request)
	fmt.Printf("Response from API is : %v\n", resp.Body)
}
