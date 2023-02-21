package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	hellopb "github.com/kakke18/grpc-practice/pkg/grpc"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func main() {
	// 8080番portのListenerを作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// gRPCサーバーを作成し、自作サービスを登録
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			myUnaryServerInterceptor1,
			myUnaryServerInterceptor2,
		),
		grpc.ChainStreamInterceptor(
			myStreamServerInterceptor1,
			myStreamServerInterceptor2,
		),
	)
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

// Hello Unary RPCがレスポンスを返す
func (s *myServer) Hello(_ context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	name := req.GetName()
	if name == "unknown" {
		// "unknown"とリクエストされたらエラーを返すようにする
		stat := status.New(codes.Unknown, "unknown error occurred")
		stat, _ = stat.WithDetails(&errdetails.DebugInfo{
			StackEntries: nil,
			Detail:       "detail reason of error",
		})

		return nil, stat.Err()
	}

	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", name),
	}, nil
}

// HelloServerStream Server Stream RPCがレスポンスを返す
func (s *myServer) HelloServerStream(req *hellopb.HelloRequest, stream hellopb.GreetingService_HelloServerStreamServer) error {
	resCount := 5
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&hellopb.HelloResponse{
			Message: fmt.Sprintf("[%d] Hello, %s!", i, req.GetName()),
		}); err != nil {
			return err
		}

		time.Sleep(time.Second * 1)
	}

	return nil
}

// HelloClientStream Client Stream RPCがリクエストを受け取る
func (s *myServer) HelloClientStream(stream hellopb.GreetingService_HelloClientStreamServer) error {
	nameList := make([]string, 0)
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			message := fmt.Sprintf("Hello, %v!", nameList)
			return stream.SendAndClose(&hellopb.HelloResponse{
				Message: message,
			})
		}

		if err != nil {
			return err
		}

		nameList = append(nameList, req.GetName())
	}
}

// HelloBiStreams 双方向ストリーミング
func (s *myServer) HelloBiStreams(stream hellopb.GreetingService_HelloBiStreamsServer) error {
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		message := fmt.Sprintf("Hello, %s!", req.GetName())
		if err := stream.Send(&hellopb.HelloResponse{
			Message: message,
		}); err != nil {
			return err
		}
	}
}
