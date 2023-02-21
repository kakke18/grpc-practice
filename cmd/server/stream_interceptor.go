package main

import (
	"errors"
	"io"
	"log"

	"google.golang.org/grpc"
)

func myStreamServerInterceptor1(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// ストリームがopenされたときに行われる前処理
	log.Println("[pre stream] my stream server interceptor 1: ", info.FullMethod)

	err := handler(srv, &myServerStreamWrapper1{ss}) // 本来の処理

	// ストリームがcloseされるときに行われる後処理
	log.Println("[post stream] my stream server interceptor 1")

	return err
}

type myServerStreamWrapper1 struct {
	grpc.ServerStream
}

func (w *myServerStreamWrapper1) RecvMsg(m interface{}) error {
	// ストリームから、リクエストを受信
	err := w.ServerStream.RecvMsg(m)
	// 受信したリクエストを、ハンドラで処理する前に差し込む前処理
	if !errors.Is(err, io.EOF) {
		log.Println("[pre message] my stream server interceptor 1: ", m)
	}
	return err
}

func (w *myServerStreamWrapper1) SendMsg(m interface{}) error {
	// ハンドラで作成したレスポンスを、ストリームから返信する直前に差し込む後処理
	log.Println("[post message] my stream server interceptor 1: ", m)
	return w.ServerStream.SendMsg(m)
}

func myStreamServerInterceptor2(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// ストリームがopenされたときに行われる前処理
	log.Println("[pre stream] my stream server interceptor 2: ", info.FullMethod)

	err := handler(srv, &myServerStreamWrapper2{ss}) // 本来の処理

	// ストリームがcloseされるときに行われる後処理
	log.Println("[post stream] my stream server interceptor 2")

	return err
}

type myServerStreamWrapper2 struct {
	grpc.ServerStream
}

func (w *myServerStreamWrapper2) RecvMsg(m interface{}) error {
	// ストリームから、リクエストを受信
	err := w.ServerStream.RecvMsg(m)
	// 受信したリクエストを、ハンドラで処理する前に差し込む前処理
	if !errors.Is(err, io.EOF) {
		log.Println("[pre message] my stream server interceptor 2: ", m)
	}
	return err
}

func (w *myServerStreamWrapper2) SendMsg(m interface{}) error {
	// ハンドラで作成したレスポンスを、ストリームから返信する直前に差し込む後処理
	log.Println("[post message] my stream server interceptor 2: ", m)
	return w.ServerStream.SendMsg(m)
}
