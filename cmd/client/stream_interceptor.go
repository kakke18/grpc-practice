package main

import (
	"context"
	"errors"
	"io"
	"log"

	"google.golang.org/grpc"
)

func myStreamClientInterceptor1(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	// ストリームがopenされたときに行われる前処理
	log.Println("[pre stream] my stream client interceptor 1: ", method)

	stream, err := streamer(ctx, desc, cc, method, opts...)

	return &myClientStreamWrapper1{stream}, err
}

type myClientStreamWrapper1 struct {
	grpc.ClientStream
}

func (w *myClientStreamWrapper1) SendMsg(m interface{}) error {
	log.Println("[pre message] my stream client interceptor 1: ", m)
	return w.ClientStream.SendMsg(m)
}

func (w *myClientStreamWrapper1) RecvMsg(m interface{}) error {
	err := w.ClientStream.RecvMsg(m)
	if !errors.Is(err, io.EOF) {
		log.Println("[post message] my stream client interceptor 1: ", m)
	}

	return err
}

func (w *myClientStreamWrapper1) CloseSend() error {
	err := w.ClientStream.CloseSend()

	log.Println("[post stream] my stream client interceptor 1")

	return err
}

func myStreamClientInterceptor2(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	// ストリームがopenされたときに行われる前処理
	log.Println("[pre stream] my stream client interceptor 2: ", method)

	stream, err := streamer(ctx, desc, cc, method, opts...)

	return &myClientStreamWrapper2{stream}, err
}

type myClientStreamWrapper2 struct {
	grpc.ClientStream
}

func (w *myClientStreamWrapper2) SendMsg(m interface{}) error {
	log.Println("[pre message] my stream client interceptor 2: ", m)
	return w.ClientStream.SendMsg(m)
}

func (w *myClientStreamWrapper2) RecvMsg(m interface{}) error {
	err := w.ClientStream.RecvMsg(m)
	if !errors.Is(err, io.EOF) {
		log.Println("[post message] my stream client interceptor 2: ", m)
	}

	return err
}

func (w *myClientStreamWrapper2) CloseSend() error {
	err := w.ClientStream.CloseSend()

	log.Println("[post stream] my stream client interceptor 2")

	return err
}
