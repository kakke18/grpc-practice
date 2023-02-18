package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	hellopb "github.com/kakke18/grpc-practice/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 8080番portのListenerを作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// gRPCサーバーを作成し、自作サービスを登録
	server := grpc.NewServer()
	hellopb.RegisterGreetingServiceServer(server, newMyServer())

	// gRPCurlを使うために、サーバリフレクションの設定
	reflection.Register(server)

	// 作成したgRPCサーバーを稼働させる
	go func() {
		log.Printf("start gRPC server port: %d\n", port)
		server.Serve(listener)
	}()

	// Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Printf("stopping gRPC server...\n")
	server.GracefulStop()
}

type myServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

func newMyServer() *myServer {
	return &myServer{}
}

func (s *myServer) Hello(_ context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	// reqからnameフィールドを取り出してresを生成
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}
