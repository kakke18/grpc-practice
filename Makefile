.PHONY: tool run-server run-client
tool:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

run-server:
	go run cmd/server/main.go cmd/server/unary_interceptor.go cmd/server/stream_interceptor.go

run-client:
	go run cmd/client/main.go cmd/client/unary_interceptor.go cmd/client/stream_interceptor.go
