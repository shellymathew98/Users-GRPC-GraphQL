package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/shellymathew98/grpc-users/users/proto"
)

var addr string = "localhost:50051"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect %v\n", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	//runCreateUser(client, "ssss", "ppppp")
	//runCreateUser(client, "aaaaa", "bbbbb")
	runGetUser(client, "2081")
	//runUpdateUser(client, "4059", "eeeee", "ffff")
	//runDeleteUser(client, "1847")

}

func runGetUser(client pb.UserServiceClient, id string) {
	log.Printf("GetUsers was invoked")
	req := &pb.Id{Value: id}
	res, err := client.GetUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	log.Printf("UserInfo: %v", res)
}

func runCreateUser(client pb.UserServiceClient, name string, place string) {
	req := &pb.UserInfo{Name: name, Place: place}
	res, err := client.CreateUser(context.Background(), req)
	if err != nil {
		log.Fatal("Errror while creating user %v\n", err)
	}
	log.Printf("Created User with Id %v\n", res)
}

func runUpdateUser(client pb.UserServiceClient, id string, name string, place string) {
	req := &pb.UserInfo{Id: id, Name: name, Place: place}
	res, err := client.UpdateUser(context.Background(), req)
	if err != nil {
		log.Fatal("Errror while updating user %v\n", err)
	}
	if res.GetValue() == 1 {
		log.Printf("Updation Success \n")
	} else {
		log.Printf("Updation Failed \n")
	}

}

func runDeleteUser(client pb.UserServiceClient, id string) {
	req := &pb.Id{Value: id}
	res, err := client.DeleteUser(context.Background(), req)
	if err != nil {
		log.Fatal("Errror while deleting user %v\n", err)
	}
	if res.GetValue() == 1 {
		log.Printf("Deletion Success \n")
	} else {
		log.Printf("Deletion Failed \n")
	}
}
