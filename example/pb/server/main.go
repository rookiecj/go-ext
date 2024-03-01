package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"protobuf-example/proto/dto"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
)

func Post(resp http.ResponseWriter, req *http.Request) {
	contentLength := req.ContentLength

	fmt.Printf("%s %s Len=%v\n", req.Method, req.URL.Path, contentLength)
	request := &dto.PostRequest{}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalf("error to read request : %v", err)
	}

	proto.Unmarshal(data, request)
	id := request.GetId()
	result := &dto.PostResponse{UserId: id, Id: id, Title: fmt.Sprintf("Title %d", id), Body: "Hello World"}
	response, err := proto.Marshal(result)
	if err != nil {
		log.Fatalf("error to marshal response : %v", err)
	}
	resp.Header().Add("Content-Type", "application/protobuf")
	resp.Write(response)
}

func main() {
	addr := "0.0.0.0:8080"
	fmt.Println("Starting the API server... ", addr)

	r := mux.NewRouter()
	r.HandleFunc("/post", Post).Methods("POST")

	server := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
