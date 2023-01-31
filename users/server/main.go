package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"strconv"

	pb "github.com/shellymathew98/grpc-users/users/proto"
	"github.com/shellymathew98/grpc-users/users/store"
	"google.golang.org/grpc"
)

var addr string = "localhost:50051"

var dbURI string = "projects/your-project-id/instances/test-instance/databases/users"

var users []*pb.UserInfo

type Server struct {
	pb.UserServiceServer
}

func main() {
	store.SpannerDbInit(dbURI)
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

	user, err := store.GetUser(in.Value, dbURI)
	if err != nil {
		return &pb.UserInfo{}, err
	}

	res := &pb.UserInfo{
		Id:    user.Id,
		Name:  user.Name,
		Place: user.Place,
	}

	return res, nil
}

func (s *Server) CreateUser(ctx context.Context, in *pb.UserInfo) (*pb.Id, error) {
	log.Printf("Recieved in %v\n", in)

	res := pb.Id{}

	newUser := store.UserInfo{
		Name:  in.Name,
		Place: in.Place,
		Id:    strconv.Itoa(rand.Intn(10000)),
	}

	res.Value = newUser.Id
	in.Id = res.GetValue()
	_, err := store.CreateUser(dbURI, newUser)
	if err != nil {
		log.Fatalf("Couldn't create user: %v", err)
	}
	users = append(users, in)
	return &res, nil
}

func (s *Server) UpdateUser(ctx context.Context, in *pb.UserInfo) (*pb.Status, error) {
	log.Printf("Recieved in %v\n", in)

	user := store.UserInfo{
		Id:    in.Id,
		Name:  in.Name,
		Place: in.Place,
	}

	err := store.UpdateUser(dbURI, user)
	if err != nil {
		return &pb.Status{Value: int32(-1)}, err
	} else {
		return &pb.Status{Value: int32(1)}, nil
	}

}

func (s *Server) DeleteUser(ctx context.Context, in *pb.Id) (*pb.Status, error) {
	log.Printf("Recieved in %v\n", in)

	err := store.DeleteUser(dbURI, in.Value)
	if err != nil {
		return &pb.Status{Value: int32(-1)}, err
	} else {
		return &pb.Status{Value: int32(1)}, nil
	}
}
