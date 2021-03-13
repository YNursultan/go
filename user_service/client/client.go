package main

import (
	"context"
	servicepb "final_project/user_service/api"
	"fmt"
	"log"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client: ")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := servicepb.NewUserServiceClient(conn)
	createUser(c)
}

func createUser(c servicepb.UserServiceClient) {
	ctx := context.Background()

	request := &servicepb.CreateUserRequest{User: &servicepb.User{
		Name:  "Azamat",
		Login: "aza123",
		Pass: "123",
		Phone: "87005488851",
	}}

	response, err := c.CreateUser(ctx, request)
	if err != nil {
		log.Fatalf("error while calling User service RPC %v", err)
	}
	log.Printf("response from User service:%v", response.Status)
}

