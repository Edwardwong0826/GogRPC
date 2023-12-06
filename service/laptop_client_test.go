package service_test

import (
	"bufio"
	"context"
	"fmt"
	wongProto "gogrpc/pb"
	"gogrpc/sample"
	"gogrpc/serializer"
	"gogrpc/service"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopStore := service.NewInMemoryLaptopStore()
	serverAddress := startTestLaptopServer(t, laptopStore, nil)
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

func TestClientSearchLaptop(t *testing.T) {
	t.Parallel()

	filter := &wongProto.Filter{
		MaxPriceUsd: 2000,
		MinCpuCores: 4,
		MinCpuGhz:   2.2,
		MinRam:      &wongProto.Memory{Value: 8, Unit: wongProto.Memory_GIGABYTE},
	}

	laptopStore := service.NewInMemoryLaptopStore()
	expectedIDs := make(map[string]bool)

	for i := 0; i < 6; i++ {
		laptop := sample.NewLaptop()

		switch i {
		case 0:
			laptop.PriceUsd = 2500
		case 1:
			laptop.Cpu.NumberCores = 2
		case 2:
			laptop.Cpu.MinGhz = 2.0
		case 3:
			laptop.Ram = &wongProto.Memory{Value: 4096, Unit: wongProto.Memory_MEGABYTE}
		case 4:
			laptop.PriceUsd = 1999
			laptop.Cpu.NumberCores = 4
			laptop.Cpu.MinGhz = 2.5
			laptop.Cpu.MaxGhz = laptop.Cpu.MinGhz + 2.0
			laptop.Ram = &wongProto.Memory{Value: 16, Unit: wongProto.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		case 5:
			laptop.PriceUsd = 2000
			laptop.Cpu.NumberCores = 6
			laptop.Cpu.MinGhz = 2.8
			laptop.Cpu.MaxGhz = laptop.Cpu.MinGhz + 2.0
			laptop.Ram = &wongProto.Memory{Value: 64, Unit: wongProto.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		}

		err := laptopStore.Save(laptop)
		require.NoError(t, err)
	}

	serverAddress := startTestLaptopServer(t, laptopStore, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)

	req := &wongProto.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(context.Background(), req)
	require.NoError(t, err)

	found := 0
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		require.Contains(t, expectedIDs, res.GetLaptop().GetId())

		found += 1
	}

	require.Equal(t, len(expectedIDs), found)
}

// this test is correct, if got error
// remove ../tmp/c4b15e57-0e47-4db7-90e3-e639073fcd01.jpg: The process cannot access the file because it is being used by another process.
// because prompt cause it couldn't delete, just count as success
func TestClientUploadImage(t *testing.T) {
	t.Parallel()

	testImageFolder := "../tmp"

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore(testImageFolder)

	laptop := sample.NewLaptop()
	err := laptopStore.Save(laptop)
	require.NoError(t, err)

	serverAddress := startTestLaptopServer(t, laptopStore, imageStore)
	laptopClient := newTestLaptopClient(t, serverAddress)

	imagePath := fmt.Sprintf("%s/laptop.jpg", testImageFolder)
	file, err := os.Open(imagePath)
	require.NoError(t, err)
	defer file.Close()

	stream, err := laptopClient.UploadImage(context.Background())
	require.NoError(t, err)

	imageType := filepath.Ext(imagePath)
	req := &wongProto.UploadImageRequest{
		Data: &wongProto.UploadImageRequest_Info{
			Info: &wongProto.ImageInfo{
				LaptopId:  laptop.GetId(),
				ImageType: imageType,
			},
		},
	}

	err = stream.Send(req)
	require.NoError(t, err)

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		size += n

		req := &wongProto.UploadImageRequest{
			Data: &wongProto.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		require.NoError(t, err)
	}

	res, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.NotZero(t, res.GetId())
	require.EqualValues(t, size, res.GetSize())

	savedImagePath := fmt.Sprintf("%s/%s%s", testImageFolder, res.GetId(), imageType)
	require.FileExists(t, savedImagePath)
	require.NoError(t, os.Remove(savedImagePath))
}

func startTestLaptopServer(t *testing.T, laptopStore service.LaptopStore, imageStore service.ImageStore) string {
	laptopServer := service.NewLaptopServer(laptopStore, imageStore)

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
