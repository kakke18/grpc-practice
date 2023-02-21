package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func myUnaryClientInterceptor1(ctx context.Context, method string, req, res interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("[pre] my unary client interceptor1", method, req)
	err := invoker(ctx, method, req, res, cc, opts...)
	fmt.Println("[post] my unary client interceptor1", res)
	return err
}

func myUnaryClientInterceptor2(ctx context.Context, method string, req, res interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("[pre] my unary client interceptor2", method, req)
	err := invoker(ctx, method, req, res, cc, opts...)
	fmt.Println("[post] my unary client interceptor2", res)
	return err
}
