package main

import (
	"context"
	servicepb "final_project/user_service/api"
	"fmt"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

//Server with embedded UnimplementedGreetServiceServer
type Server struct {
	servicepb.UnimplementedUserServiceServer
}

func (s *Server) CreateUser(ctx context.Context, req *servicepb.CreateUserRequest) (*servicepb.CreateUserResponse, error) {
	fmt.Printf("Create user function was invoked with %v \n", req)
	name := req.GetUser().GetName()
	login := req.GetUser().GetLogin()
	pass := req.GetUser().GetPass()
	phone := req.GetUser().GetPhone()

	status := ""

	//CREATE USER FUNC
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:Azamat2341!@localhost:5432/user_service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	conn.Exec(context.Background(), "insert into users(name, login, pass, phone) values($1, $2, $3, $4)", name, login, pass, phone)
	status = "Created user with login: "+login

	res := &servicepb.CreateUserResponse{
		Status: status,
	}

	return res, nil
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	servicepb.RegisterUserServiceServer(s, &Server{})
	log.Println("Server is running on port:50051")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}

}
