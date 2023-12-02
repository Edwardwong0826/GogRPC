package service

import (
	"context"
	"errors"
	wongProto "gogrpc/pb"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// the server that provides laptop services
type LaptopServer struct {
	Store LaptopStore
}

// returns a new LaptopServer
func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{store}
}

// CreateLaptop is a unary RPC to create a new laptop
func (server *LaptopServer) CreateLaptop(ctx context.Context, req *wongProto.CreateLaptopRequest) (*wongProto.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("receive a create-laptop request with id: %s", laptop.Id)

	if len(laptop.Id) >= 0 {
		// check if it's a valid UUID
		//go get github.com/google/uuid

		_, err := uuid.Parse(laptop.Id)

		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID: %v", err)
		} else {
			id, err := uuid.NewRandom()
			if err != nil {
				return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
			}

			laptop.Id = id.String()

		}
	}

	// save the laptop to db if want to, here just save laptop to in-memory store
	err := server.Store.Save(laptop)

	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}

		return nil, status.Errorf(code, "cannot save laptop to the store %v", err)

	}

	log.Printf("saved laptop with id %s", laptop.Id)
	res := &wongProto.CreateLaptopResponse{
		Id: laptop.Id,
	}

	return res, nil
}
