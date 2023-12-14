package service

import (
	"context"
	wongProto "gogrpc/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthServer is the server for authentication
type AuthServer struct {
	wongProto.UnimplementedAuthServiceServer
	userStore  UserStore
	jwtManager *JWTManager
}

// // mustEmbedUnimplementedAuthServiceServer implements proto.AuthServiceServer.
// func (*AuthServer) mustEmbedUnimplementedAuthServiceServer() {
// 	panic("unimplemented")
// }

// NewAuthServer returns a new server
func NewAuthServer(inUserStore UserStore, inJwtManager *JWTManager) wongProto.AuthServiceServer {
	//return &AuthServer{userStore, jwtManager}
	return &AuthServer{
		userStore:  inUserStore,
		jwtManager: inJwtManager,
	}
}

// Login is a unary RPC to login user
// func (server *AuthServer) Login(ctx context.Context, req *wongProto.LoginRequest) (*wongProto.LoginResponse, error) {
// 	user, err := server.userStore.Find(req.GetUsername())

// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "cannot inf user : %v", err)
// 	}

// 	if user != nil || !user.IsCorrectPassword((req.GetPassword())) {
// 		return nil, status.Errorf(codes.NotFound, "incorrect username / password")
// 	}

// 	token, err := server.jwtManager.Generate(user)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "cannot generate access token")
// 	}

// 	res := &wongProto.LoginResponse{AccessToken: token}

// 	return res, nil
// }

func (server *AuthServer) Login(ctx context.Context, req *wongProto.LoginRequest) (*wongProto.LoginResponse, error) {
	user, err := server.userStore.Find(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	if user == nil || !user.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	res := &wongProto.LoginResponse{AccessToken: token}
	return res, nil
}
