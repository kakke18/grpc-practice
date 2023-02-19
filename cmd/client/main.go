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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		grpc.WithTransportCredentials(insecure.NewCredentials()), // コネクションでSSL/TLSを使用しない
		grpc.WithBlock(),                                         // コネクションが確立されるまで待機する（同期処理をする）
	)
	if err != nil {
		log.Fatalf("Connection failed. err:%s\n", err)
	}
	defer conn.Close()

	// gPRCクライアントを作成
	client := hellopb.NewGreetingServiceClient(conn)

	// 入力待ち状態にする
	for {
		fmt.Printf("1: send Request\n")
		fmt.Printf("2: HelloServerStream\n")
		fmt.Printf("3: HelloClientStream\n")
		fmt.Printf("4: exit\n")
		fmt.Printf("pleace enter >")

		input := getInputString(scanner)

		switch input {
		case "1":
			res, err := hello(ctx, scanner, client)
			if err != nil {
				fmt.Printf("err:%s\n", err)
			}

			fmt.Printf("res:%s\n", res)

			return
		case "2":
			res, err := helloServerStream(ctx, scanner, client)
			if err != nil {
				fmt.Printf("err:%s\n", err)
			}

			fmt.Printf("res:%+v\n", res)

			return
		case "3":
			res, err := helloClientStream(ctx, scanner, client)
			if err != nil {
				fmt.Printf("err:%s\n", err)
			}

			fmt.Printf("res:%s\n", res)

			return
		case "4":
			fmt.Printf("bye.\n")
			return
		default:
			fmt.Printf("unexpected input. bye.\n")
			return
		}
	}
}

func hello(ctx context.Context, scanner *bufio.Scanner, client hellopb.GreetingServiceClient) (string, error) {
	fmt.Printf("Please enter your name.\n")

	name := getInputString(scanner)
	req := &hellopb.HelloRequest{
		Name: name,
	}
	res, err := client.Hello(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetMessage(), nil
}

func helloServerStream(ctx context.Context, scanner *bufio.Scanner, client hellopb.GreetingServiceClient) ([]string, error) {
	fmt.Println("Please enter your name.")

	name := getInputString(scanner)
	req := &hellopb.HelloRequest{
		Name: name,
	}
	stream, err := client.HelloServerStream(ctx, req)
	if err != nil {
		return nil, err
	}

	var res []string
	for {
		helloRes, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("all the responses have already received.")
			break
		}

		if err != nil {
			return nil, err
		}

		res = append(res, helloRes.GetMessage())
	}

	return res, nil
}

func helloClientStream(ctx context.Context, scanner *bufio.Scanner, client hellopb.GreetingServiceClient) (string, error) {
	stream, err := client.HelloClientStream(ctx)
	if err != nil {
		return "", err
	}

	seedCount := 5
	fmt.Printf("Please enter %d names.\n", seedCount)
	for i := 0; i < seedCount; i++ {
		name := getInputString(scanner)
		if err := stream.Send(&hellopb.HelloRequest{
			Name: name,
		}); err != nil {
			return "", err
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}

	return res.GetMessage(), nil
}

func getInputString(scanner *bufio.Scanner) string {
	scanner.Scan()
	in := scanner.Text()

	return in
}
