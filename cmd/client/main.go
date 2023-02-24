package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	hellopb "github.com/kakke18/grpc-practice/pkg/grpc"
	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	ctx := context.Background()

	fmt.Printf("start gRPC Client.\n")

	// 標準入力から文字列を受け取るスキャナを用意
	scanner := bufio.NewScanner(os.Stdin)

	// コネクションを確立
	address := "localhost:8080"
	conn, err := grpc.Dial(
		address,
		grpc.WithChainUnaryInterceptor(
			myUnaryClientInterceptor1,
			myUnaryClientInterceptor2,
		),
		grpc.WithChainStreamInterceptor(
			myStreamClientInterceptor1,
			myStreamClientInterceptor2,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()), // コネクションでSSL/TLSを使用しない
		grpc.WithBlock(), // コネクションが確立されるまで待機する（同期処理をする）
	)
	if err != nil {
		log.Fatalf("Connection failed. err:%s\n", err)
	}
	defer conn.Close()

	// gPRCクライアントを作成
	client := hellopb.NewGreetingServiceClient(conn)

	// 入力待ち状態にする
	for {
		fmt.Println("1    : send Request")
		fmt.Println("2    : HelloServerStream")
		fmt.Println("3    : HelloClientStream")
		fmt.Println("4    : HelloBiStream")
		fmt.Println("other: exit")
		fmt.Printf("pleace enter > ")

		input := getInputString(scanner)

		switch input {
		case "1":
			hello(ctx, scanner, client)
		case "2":
			helloServerStream(ctx, scanner, client)
		case "3":
			helloClientStream(ctx, scanner, client)
		case "4":
			helloBiStream(ctx, scanner, client)
		default:
			fmt.Println("bye.")
		}

		return
	}
}

func hello(ctx context.Context, scanner *bufio.Scanner, client hellopb.GreetingServiceClient) {
	fmt.Printf("please enter your name > ")

	name := getInputString(scanner)
	req := &hellopb.HelloRequest{
		Name: name,
	}
	md := metadata.New(map[string]string{
		"type": "unary",
		"from": "client",
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	var header, trailer metadata.MD
	res, err := client.Hello(ctx, req, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		if stat, ok := status.FromError(err); ok {
			fmt.Printf("code: %s, message: %s, details: %s\n", stat.Code(), stat.Message(), stat.Details())
		} else {
			fmt.Printf("err: %s\n", err.Error())
		}
		return
	}

	fmt.Printf("header: %+v, trailer: %+v, res: %s\n", header, trailer, res.GetMessage())
}

func helloServerStream(ctx context.Context, scanner *bufio.Scanner, client hellopb.GreetingServiceClient) {
	fmt.Printf("please enter your name > ")

	name := getInputString(scanner)
	req := &hellopb.HelloRequest{
		Name: name,
	}
	stream, err := client.HelloServerStream(ctx, req)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("all the responses have already received.")
			} else {
				fmt.Printf("err: %s\n", err.Error())
			}
			return
		}

		fmt.Printf("res: %s\n", res.GetMessage())
	}
}

func helloClientStream(ctx context.Context, scanner *bufio.Scanner, client hellopb.GreetingServiceClient) {
	stream, err := client.HelloClientStream(ctx)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return
	}

	seedCount := 5
	fmt.Printf("Please enter %d names.\n", seedCount)
	for i := 0; i < seedCount; i++ {
		name := getInputString(scanner)
		if err := stream.Send(&hellopb.HelloRequest{
			Name: name,
		}); err != nil {
			fmt.Printf("err: %s\n", err.Error())
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return
	}

	fmt.Printf("res: %s\n", res.GetMessage())
}

func helloBiStream(ctx context.Context, scanner *bufio.Scanner, client hellopb.GreetingServiceClient) {
	md := metadata.New(map[string]string{
		"type": "stream",
		"from": "client",
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.HelloBiStreams(ctx)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return
	}

	seedCount := 5
	fmt.Printf("Please enter %d names.\n", seedCount)

	var (
		sendEnd, recvEnd bool
		sendCount        int
	)
	for !(sendEnd && recvEnd) {
		// 送信処理
		if !sendEnd {
			name := getInputString(scanner)
			sendCount++
			if err := stream.Send(&hellopb.HelloRequest{
				Name: name,
			}); err != nil {
				fmt.Printf("err: %s\n", err.Error())
				return
			}

			if sendCount == seedCount {
				sendEnd = true
				if err := stream.CloseSend(); err != nil {
					fmt.Printf("err: %s\n", err.Error())
					return
				}
			}
		}

		// 受信処理
		var (
			headerMD metadata.MD
			helloRes *hellopb.HelloResponse
		)
		if !recvEnd {
			if headerMD == nil {
				headerMD, err = stream.Header()
				if err != nil {
					fmt.Printf("err: %s\n", err.Error())
				} else {
					fmt.Printf("headerMD: %+v\n", headerMD)
				}
			}

			res, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Println("all the responses have already received.")
					break
				} else {
					fmt.Printf("err: %s\n", err.Error())
				}
				return
			}

			helloRes = res
		}

		fmt.Printf("res: %s\n", helloRes.GetMessage())
	}

	trailerMD := stream.Trailer()
	fmt.Printf("trailerMD: %+v\n", trailerMD)
}

func getInputString(scanner *bufio.Scanner) string {
	scanner.Scan()
	in := scanner.Text()

	return in
}
