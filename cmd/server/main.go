package main

import (
	"flag"
	"fmt"
	"gogrpc/service"
	"log"
	"net"
	"time"

	wongProto "gogrpc/pb"

	"google.golang.org/grpc"
)

func seedUsers(userStore service.UserStore) error {
	err := createUser(userStore, "admin1", "secret", "admin")
	if err != nil {
		return err
	}
	return createUser(userStore, "user1", "secret", "user")
}

func createUser(userStore service.UserStore, username, password, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}
	return userStore.Save(user)
}

// // function name define ourself, just method parameter follow to intercetpor.go and need to add grpc into the fields
// // type UnaryServerInterceptor func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (resp any, err error)
// func UnaryInterceptor(
// 	ctx context.Context,
// 	req any,
// 	info *grpc.UnaryServerInfo,
// 	handler grpc.UnaryHandler,
// ) (resp any, err error) {
// 	log.Println("--> unary interceptor: ", info.FullMethod)
// 	return handler(ctx, req)
// }

// // function name define ourself, just method parameter follow to intercetpor.go and need to add grpc into the fields
// // type StreamServerInterceptor func(srv any, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error
// func streamInterceptor(
// 	srv any,
// 	ss grpc.ServerStream,
// 	info *grpc.StreamServerInfo,
// 	handler grpc.StreamHandler) error {
// 	log.Println("--> stream interceptor: ", info.FullMethod)
// 	return handler(srv, ss)
// }

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

func accessibleRoles() map[string][]string {
	const laptopServicePath = "/wong.LaptopService/"

	return map[string][]string{
		laptopServicePath + "CreateLaptop": {"admin"},
		laptopServicePath + "UploadImage":  {"admin"},
		laptopServicePath + "RateLaptop":   {"admin", "user"},
	}
}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port %d", *port)

	userStore := service.NewInMemoryUserStore()
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users: ", err)
	}

	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(userStore, jwtManager)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	ratingStore := service.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)

	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())

	// There are two types of interceptor in gRPC, one for unary RPC and the other for stream RPC
	grpcServer := grpc.NewServer(
		// pass in method parameter need to go to this grpc.UnaryInterceptor and look the definition in interceptor.go
		// type UnaryServerInterceptor func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (resp any, err error)

		//grpc.UnaryInterceptor(UnaryInterceptor),
		//grpc.StreamInterceptor(streamInterceptor),
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)

	wongProto.RegisterAuthServiceServer(grpcServer, authServer)
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
