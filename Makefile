.PHONY: tool run-server run-client
tool:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go
