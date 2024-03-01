
.PHONY: setup
setup:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go get google.golang.org/protobuf/proto
	go get github.com/gorilla/mux

dep:
	go get -u ./...
	go mod tidy

.PHONY: proto
proto:
	protoc -I=./proto --go_out=./proto ./proto/post.proto


.PHONY: clean
clean:
	-rm -rf ./proto/dto

.PHONY: server
server:
	-go run server/main.go

.PHONY: client
client: server
	go run cmd/httpx_bodyparser.go
	