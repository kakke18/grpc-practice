package main

import (
	"bufio"
	"context"
	"fmt"
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
		fmt.Printf("2: exit\n")
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

func getInputString(scanner *bufio.Scanner) string {
	scanner.Scan()
	in := scanner.Text()

	return in
}
