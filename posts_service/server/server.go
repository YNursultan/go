package main

import (
	"context"
	servicepb "final_project/posts_service/api"
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
	servicepb.UnimplementedPostServiceServer
}

func (s *Server) CreatePost(ctx context.Context, req *servicepb.CreatePostRequest) (*servicepb.CreatePostResponse, error) {
	fmt.Printf("Create post function was invoked with %v \n", req)
	title := req.GetPost().GetTitle()
	desc := req.GetPost().GetDescription()
	category := req.GetPost().GetCategory()
	user_id := req.GetPost().GetUserId()

	status := ""

	//CREATE POST FUNC
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:1234@localhost:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	conn.Exec(context.Background(), "insert into posts(title, description, category, user_id) values($1, $2, $3, $4)",
		title, desc, category, user_id)
	status = "Created post with title: " + title

	res := &servicepb.CreatePostResponse{
		Status: status,
	}

	return res, nil
}

func (s *Server) UpdatePost(ctx context.Context, req *servicepb.UpdatePostRequest) (*servicepb.UpdatePostResponse, error) {
	fmt.Printf("Update post function was invoked with %v \n", req)
	id := req.GetPost().GetId()
	title := req.GetPost().GetTitle()
	desc := req.GetPost().GetDescription()
	category := req.GetPost().GetCategory()
	user_id := req.GetPost().GetUserId()

	status := ""

	//UPDATE POST FUNC
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:1234@localhost:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	conn.Exec(context.Background(), "update posts set title=$1, description=$2, category=$3, user_id=$4 where id = $5",
		title, desc, category, user_id, id)
	status = "Updated post with title: " + title

	res := &servicepb.UpdatePostResponse{
		Status: status,
	}

	return res, nil
}

func (s *Server) DeletePost(ctx context.Context, req *servicepb.DeletePostRequest) (*servicepb.DeletePostResponse, error) {
	fmt.Printf("Delete post function was invoked with %v \n", req)
	id := req.PostId
	status := ""

	//DELETED USER FUNC
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:1234@localhost:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	conn.Exec(context.Background(), "delete from posts where id = $1", id)
	status = "Deleted post with id: " + strconv.Itoa(int(id))

	res := &servicepb.DeletePostResponse{
		Status: status,
	}

	return res, nil
}

func (s *Server) GetPost(ctx context.Context, req *servicepb.GetPostRequest) (*servicepb.GetPostResponse, error) {
	fmt.Printf("Get user function was invoked with %v \n", req)
	id := req.PostId
	title := ""
	desc := ""
	category := ""
	user_id := 0

	//GET POST FUNC
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:1234@localhost:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	err = conn.QueryRow(context.Background(), "select * from posts where id = $1", id).Scan(&id, &title, &desc, &category, &user_id)

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	res := &servicepb.GetPostResponse{
		Post: &servicepb.Post{Id: int64(id), Title: title, Description: desc, Category: category, UserId: int64(user_id)},
	}

	return res, nil
}

func (s *Server) GetAllPosts(req *servicepb.GetAllPostsRequest, stream servicepb.PostService_GetAllPostsServer) error {
	fmt.Printf("Get all posts function was invoked with %v \n", req)

	stml := "SELECT * FROM POSTS"

	conn, err := pgx.Connect(context.Background(), "postgres://postgres:1234@localhost:5432/postgres")

	rows, err := conn.Query(context.Background(), stml)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		s := &servicepb.Post{}

		err = rows.Scan(&s.Id, &s.Title, &s.Description, &s.Category, &s.UserId)
		if err != nil {
			fmt.Printf("Error: %s", err)
		}

		res := &servicepb.GetAllPostsResponse{Post: s}
		if err := stream.Send(res); err != nil {
			log.Fatalf("error while sending posts responses: %v", err.Error())
		}
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return nil
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	servicepb.RegisterPostServiceServer(s, &Server{})
	log.Println("Server is running on port:50051")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}

}
