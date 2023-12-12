package main

import (
	"context"
	"flag"
	"fmt"
	"gogrpc/service"
	"log"
	"net"

	wongProto "gogrpc/pb"

	"google.golang.org/grpc"
)

// function name define ourself, just method parameter follow to intercetpor.go and need to add grpc into the fields
// type UnaryServerInterceptor func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (resp any, err error)
func UnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	log.Println("--> unary interceptor: ", info.FullMethod)
	return handler(ctx, req)
}

// function name define ourself, just method parameter follow to intercetpor.go and need to add grpc into the fields
// type StreamServerInterceptor func(srv any, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error
func streamInterceptor(
	srv any,
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	log.Println("--> stream interceptor: ", info.FullMethod)
	return handler(srv, ss)
}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port %d", *port)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	ratingStore := service.NewInMemoryRatingStore()

	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)
	// There are two types of interceptor in gRPC, one for unary RPC and the other for stream RPC
	grpcServer := grpc.NewServer(
		// pass in method parameter need to go to this grpc.UnaryInterceptor and look the definition in interceptor.go
		// type UnaryServerInterceptor func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (resp any, err error)
		grpc.UnaryInterceptor(UnaryInterceptor),
		grpc.StreamInterceptor(streamInterceptor),
	)
	wongProto.RegisterLaptopServiceServer(grpcServer, laptopServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server: ", err)
	}

}
