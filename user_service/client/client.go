package main

import (
	"context"
	servicepb "final_project/user_service/api"
	"fmt"
	"google.golang.org/grpc"
	"log"
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
		Pass:  "123",
		Phone: "87005488851",
	}}

	response, err := c.CreateUser(ctx, request)
	if err != nil {
		log.Fatalf("error while calling User server RPC %v", err)
	}
	log.Printf("response from User server:%v", response.Status)
}

func getUser(c servicepb.UserServiceClient) {
	ctx := context.Background()

	request := &servicepb.GetUserRequest{UserId: 3}

	response, err := c.GetUser(ctx, request)
	if err != nil {
		log.Fatalf("error while calling User server RPC %v", err)
	}
	log.Printf("response from User server:%v", response.User.Login)
}

func deleteUser(c servicepb.UserServiceClient) {
	ctx := context.Background()

	request := &servicepb.DeleteUserRequest{UserId: 3}

	response, err := c.DeleteUser(ctx, request)
	if err != nil {
		log.Fatalf("error while calling User server RPC %v", err)
	}
	log.Printf("response from User server:%v", response.Status)
}

func updateUser(c servicepb.UserServiceClient) {
	ctx := context.Background()

	request := &servicepb.UpdateUserRequest{User: &servicepb.User{Login: "SS", Name: "AZA", Pass: "1234", Phone: "123456", Id: 6}}

	response, err := c.UpdateUser(ctx, request)
	if err != nil {
		log.Fatalf("error while calling User server RPC %v", err)
	}
	log.Printf("response from User server:%v", response.Status)
}
