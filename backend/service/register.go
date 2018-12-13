package service

import (
	"backend/pb/user"
	"backend/service/user"
	"google.golang.org/grpc"
)
// RegisterService register service
func RegisterService(server *grpc.Server)  {

	pb.RegisterUserServiceServer(server,&user.User{})
}