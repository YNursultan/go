package main

import (
	"context"
	servicepb "final_project/posts_service/api"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	fmt.Println("Client: ")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := servicepb.NewPostServiceClient(conn)
	getAllPosts(c)
}

func createPost(c servicepb.PostServiceClient) {
	ctx := context.Background()

	request := &servicepb.CreatePostRequest{Post: &servicepb.Post{
		Title:       "Iphone X",
		Description: "128GB Silver",
		Category:    "Tech",
		UserId:      6,
	}}

	response, err := c.CreatePost(ctx, request)
	if err != nil {
		log.Fatalf("error while calling Post server RPC %v", err)
	}
	log.Printf("response from Post server:%v", response.Status)
}

func getPost(c servicepb.PostServiceClient) {
	ctx := context.Background()

	request := &servicepb.GetPostRequest{PostId: 1}

	response, err := c.GetPost(ctx, request)
	if err != nil {
		log.Fatalf("error while calling Post server RPC %v", err)
	}
	log.Printf("response from Post server:%v", response.Post.Title)
}

func deletePost(c servicepb.PostServiceClient) {
	ctx := context.Background()

	request := &servicepb.DeletePostRequest{PostId: 1}

	response, err := c.DeletePost(ctx, request)
	if err != nil {
		log.Fatalf("error while calling Post server RPC %v", err)
	}
	log.Printf("response from Post server:%v", response.Status)
}

func updatePost(c servicepb.PostServiceClient) {
	ctx := context.Background()

	request := &servicepb.UpdatePostRequest{Post: &servicepb.Post{Id: 1, Title: "changed", Description: "changed_desc", Category: "Tech", UserId: 6}}

	response, err := c.UpdatePost(ctx, request)
	if err != nil {
		log.Fatalf("error while calling Post server RPC %v", err)
	}
	log.Printf("response from Post server:%v", response.Status)
}

func getAllPosts(c servicepb.PostServiceClient) {
	ctx := context.Background()
	req := &servicepb.GetAllPostsRequest{}

	stream, err := c.GetAllPosts(ctx, req)
	if err != nil {
		log.Fatalf("error while calling GET ALL POSTS RPC %v", err)
	}
	defer stream.CloseSend()

LOOP:
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break LOOP
			}
			log.Fatalf("error while reciving from get all posts RPC %v", err)
		}
		log.Printf("response from get all posts:%v \n", res.GetPost().Title)
	}
}
