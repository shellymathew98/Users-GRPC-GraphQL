package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"strconv"

	pb "github.com/shellymathew98/grpc-users/users/proto"
	"google.golang.org/grpc"
)

var addr string = "localhost:50051"

var users []*pb.UserInfo

type Server struct {
	pb.UserServiceServer
}

func main() {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to listen on: %v\n", err)
	}

	log.Printf("Listening on %s\n", addr)

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &Server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}

}

func (s *Server) GetUser(ctx context.Context, in *pb.Id) (*pb.UserInfo, error) {
	log.Printf("Recieved in %v\n", in)

	res := &pb.UserInfo{}

	for _, user := range users {
		if user.GetId() == in.GetValue() {
			res = user
			break
		}
	}
	return res, nil
}

func (s *Server) CreateUser(ctx context.Context, in *pb.UserInfo) (*pb.Id, error) {
	log.Printf("Recieved in %v\n", in)

	res := pb.Id{}

	res.Value = strconv.Itoa(rand.Intn(10000))
	in.Id = res.GetValue()
	users = append(users, in)
	return &res, nil
}

func (s *Server) UpdateUser(ctx context.Context, in *pb.UserInfo) (*pb.Status, error) {
	log.Printf("Recieved in %v\n", in)

	res := pb.Status{}
	for index, user := range users {
		if user.GetId() == in.GetId() {
			users = append(users[:index], users[index+1:]...)
			in.Id = user.GetId()
			users = append(users, in)
			res.Value = 1
			break
		}
	}
	return &res, nil
}

func (s *Server) DeleteUser(ctx context.Context, in *pb.Id) (*pb.Status, error) {
	log.Printf("Recieved in %v\n", in)

	res := pb.Status{}

	for index, user := range users {
		if user.GetId() == in.GetValue() {
			users = append(users[:index], users[index+1:]...)
			res.Value = 1
			break
		}
	}
	return &res, nil
}
