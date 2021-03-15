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
	"strconv"
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
	status = "Created user with login: " + login

	res := &servicepb.CreateUserResponse{
		Status: status,
	}

	return res, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *servicepb.UpdateUserRequest) (*servicepb.UpdateUserResponse, error) {
	fmt.Printf("Update user function was invoked with %v \n", req)
	id := req.GetUser().GetId()
	name := req.GetUser().GetName()
	login := req.GetUser().GetLogin()
	pass := req.GetUser().GetPass()
	phone := req.GetUser().GetPhone()

	status := ""

	//UPDATE USER FUNC
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:Azamat2341!@localhost:5432/user_service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	conn.Exec(context.Background(), "update users set name=$1, login=$2, pass=$3, phone=$4 where id = $5", name, login, pass, phone, id)
	status = "Updated user with login: " + login

	res := &servicepb.UpdateUserResponse{
		Status: status,
	}

	return res, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *servicepb.DeleteUserRequest) (*servicepb.DeleteUserResponse, error) {
	fmt.Printf("Delete user function was invoked with %v \n", req)
	id := req.UserId
	status := ""

	//DELETED USER FUNC
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:Azamat2341!@localhost:5432/user_service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	conn.Exec(context.Background(), "delete from users where id = $1", id)
	status = "Deleted user with id: " + strconv.Itoa(int(id))

	res := &servicepb.DeleteUserResponse{
		Status: status,
	}

	return res, nil
}

func (s *Server) GetUser(ctx context.Context, req *servicepb.GetUserRequest) (*servicepb.GetUserResponse, error) {
	fmt.Printf("Get user function was invoked with %v \n", req)
	id := req.UserId
	name := ""
	login := ""
	pass := ""
	phone := ""

	//GET USER FUNC
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:Azamat2341!@localhost:5432/user_service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	err = conn.QueryRow(context.Background(), "select * from users where id = $1", id).Scan(&id, &name, &phone, &login, &pass)

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	res := &servicepb.GetUserResponse{
		User: &servicepb.User{Id: int64(id), Name: name, Phone: phone, Pass: pass, Login: login},
	}

	return res, nil
}

func (s *Server) GetUserByLogin(ctx context.Context, req *servicepb.GetUserByLoginRequest) (*servicepb.GetUserByLoginResponse, error) {
	fmt.Printf("Get user by login function was invoked with %v \n", req)
	id := 0
	name := ""
	login := req.UserLogin
	pass := ""
	phone := ""

	//GET USER BY LOGIN FUNC
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:Azamat2341!@localhost:5432/user_service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	err = conn.QueryRow(context.Background(), "select * from users where login = $1", login).Scan(&id, &name, &phone, &login, &pass)

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	res := &servicepb.GetUserByLoginResponse{
		User: &servicepb.User{Id: int64(id), Name: name, Phone: phone, Pass: pass, Login: login},
	}

	return res, nil
}

func (s *Server) LoginUser(ctx context.Context, req *servicepb.LoginRequest) (*servicepb.LoginResponse, error) {
	fmt.Printf("Login user function was invoked with %v \n", req)
	id := 0
	name := ""
	login := ""
	pass := ""
	phone := ""

	u_login := req.Login
	u_pass := req.Pass

	//GET USER FUNC
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:Azamat2341!@localhost:5432/user_service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	err = conn.QueryRow(context.Background(), "select * from users where login = $1", u_login).Scan(&id, &name, &phone, &login, &pass)

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	status := ""

	if err == pgx.ErrNoRows {
		status = "User with this login not registered!"
		return &servicepb.LoginResponse{Status: status}, nil
	}
	if u_pass == pass {
		status = "OK"
	} else {
		status = "Wrong pass!"
	}

	return &servicepb.LoginResponse{Status: status}, nil
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
