package service_test

import (
	"context"
	wongProto "gogrpc/pb"
	"gogrpc/sample"
	"gogrpc/serializer"
	"gogrpc/service"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopStore := service.NewInMemoryLaptopStore()
	serverAddress := startTestLaptopServer(t, laptopStore)
	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedID := laptop.Id
	req := &wongProto.CreateLaptopRequest{
		Laptop: laptop,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedID, res.Id)

	// check that the laptop is saved to the store
	other, err := laptopStore.Find(res.Id)
	require.NoError(t, err)
	require.NotNil(t, other)

	// check that the saved laptop is the same as the one we send
	requireSameLaptop(t, laptop, other)
}

func startTestLaptopServer(t *testing.T, laptopStore service.LaptopStore) string {
	laptopServer := service.NewLaptopServer(laptopStore)

	grpcServer := grpc.NewServer()
	wongProto.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0") // random available port
	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return listener.Addr().String()
}

func newTestLaptopClient(t *testing.T, serverAddress string) wongProto.LaptopServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	require.NoError(t, err)
	return wongProto.NewLaptopServiceClient(conn)
}

func requireSameLaptop(t *testing.T, laptop1 *wongProto.Laptop, laptop2 *wongProto.Laptop) {
	json1, err := serializer.ProtobufToJSON(laptop1)
	require.NoError(t, err)

	json2, err := serializer.ProtobufToJSON(laptop2)
	require.NoError(t, err)

	require.Equal(t, json1, json2)
}

// func TestClientCreateLaptop(t *testing.T) {
// 	t.Parallel()

// 	laptopStore := service.NewInMemoryLaptopStore()
// 	_, serverAddress := startTestLaptopServer(t)
// 	laptopClient := newTestLaptopClient(t, serverAddress)

// 	laptop := sample.NewLaptop()
// 	expectedID := laptop.Id
// 	req := &wongProto.CreateLaptopRequest{
// 		Laptop: laptop,
// 	}

// 	res, err := laptopClient.CreateLaptop(context.Background(), req)
// 	require.NoError(t, err)
// 	require.NotNil(t, res)
// 	require.Equal(t, expectedID, res.Id)

// 	// check that the laptop is saved to the store
// 	other, err := laptopStore.Find(res.Id)
// 	require.NoError(t, err)
// 	require.NotNil(t, other)

// 	// check that the saved laptop is the same as the one we send
// 	requireSameLaptop(t, laptop, other)
// }

// // return string is network address of the server
// func startTestLaptopServer(t *testing.T) (*service.LaptopServer, string) {
// 	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore())

// 	grpcServer := grpc.NewServer()
// 	wongProto.RegisterLaptopServiceServer(grpcServer, laptopServer)

// 	listener, err := net.Listen("tcp", ":0") // random available port
// 	require.NoError(t, err)

// 	go grpcServer.Serve(listener) // block call

// 	return laptopServer, listener.Addr().String()

// 	// go get google.golang.org/grpc/internal/transport@v1.59.0
// 	// go get google.golang.org/grpc@v1.59.0
// 	// if after isntall package got error importing for missing xxx entry for module providing package xxx; to add: go get xxx
// 	// run go mod tidy

// }

// func newTestLaptopClient(t *testing.T, serverAddress string) wongProto.LaptopServiceClient {
// 	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
// 	require.NoError(t, err)
// 	return wongProto.NewLaptopServiceClient(conn)

// }

// func requireSameLaptop(t *testing.T, laptop1 *wongProto.Laptop, laptop2 *wongProto.Laptop) {
// 	json1, err := serializer.ProtobufToJSON(laptop1)
// 	require.NoError(t, err)

// 	json2, err := serializer.ProtobufToJSON(laptop2)
// 	require.NoError(t, err)

// 	require.Equal(t, json1, json2)
// }
