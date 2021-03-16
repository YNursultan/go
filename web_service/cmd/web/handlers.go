package main

import (
	"context"
	servicepb1 "final_project/posts_service/api"
	servicepb2 "final_project/user_service/api"
	"final_project/web_service/pkg/models"
	"github.com/jackc/pgx/v4"
	"io"
	"log"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.getAllPosts(app.c1, w, r)
}

// User service adapter
func (app *application) createUserAdapter(w http.ResponseWriter, r *http.Request) {
	createUser(app.c2)
}
func (app *application) getUserAdapter(w http.ResponseWriter, r *http.Request) {
	getUser(app.c2)
}
func (app *application) deleteUserAdapter(w http.ResponseWriter, r *http.Request) {
	deleteUser(app.c2)
}
func (app *application) updateUserAdapter(w http.ResponseWriter, r *http.Request) {
	updateUser(app.c2)
}

// Post service adapter
func (app *application) createPostAdapter(w http.ResponseWriter, r *http.Request) {
	createPost(app.c1)
}
func (app *application) getPostAdapter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	app.getPost(app.c1, int32(id), w, r)
} //
func (app *application) updatePostAdapter(w http.ResponseWriter, r *http.Request) {
	updatePost(app.c1)
}
func (app *application) deletePostAdapter(w http.ResponseWriter, r *http.Request) {
	deletePost(app.c1)
}

// User service
func createUser(c servicepb2.UserServiceClient) {
	ctx := context.Background()

	request := &servicepb2.CreateUserRequest{User: &servicepb2.User{
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
func getUser(c servicepb2.UserServiceClient) {
	ctx := context.Background()

	request := &servicepb2.GetUserRequest{UserId: 3}

	response, err := c.GetUser(ctx, request)
	if err != nil {
		log.Fatalf("error while calling User server RPC %v", err)
	}
	log.Printf("response from User server:%v", response.User.Login)
}
func deleteUser(c servicepb2.UserServiceClient) {
	ctx := context.Background()

	request := &servicepb2.DeleteUserRequest{UserId: 3}

	response, err := c.DeleteUser(ctx, request)
	if err != nil {
		log.Fatalf("error while calling User server RPC %v", err)
	}
	log.Printf("response from User server:%v", response.Status)
}
func updateUser(c servicepb2.UserServiceClient) {
	ctx := context.Background()

	request := &servicepb2.UpdateUserRequest{User: &servicepb2.User{Login: "SS", Name: "AZA", Pass: "1234", Phone: "123456", Id: 6}}

	response, err := c.UpdateUser(ctx, request)
	if err != nil {
		log.Fatalf("error while calling User server RPC %v", err)
	}
	log.Printf("response from User server:%v", response.Status)
}
func loginUser(c servicepb2.UserServiceClient) {
	ctx := context.Background()

	request := &servicepb2.LoginRequest{Login: "SSs", Pass: "1234"}

	response, err := c.LoginUser(ctx, request)
	if err != nil && err != pgx.ErrNoRows {
		log.Fatalf("error while calling User LOGIN RPC %v", err)
	}
	log.Printf("response from User server:%v", response.Status)
}
func getUserByLogin(c servicepb2.UserServiceClient) {
	ctx := context.Background()

	request := &servicepb2.GetUserByLoginRequest{UserLogin: "SS"}

	response, err := c.GetUserByLogin(ctx, request)
	if err != nil {
		log.Fatalf("error while calling User server RPC %v", err)
	}
	log.Printf("response from User server:%v", response.User.Login)
}

// Post service
func createPost(c servicepb1.PostServiceClient) {
	ctx := context.Background()

	request := &servicepb1.CreatePostRequest{Post: &servicepb1.Post{
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
func (app *application) getPost(c servicepb1.PostServiceClient, id int32, w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	request := &servicepb1.GetPostRequest{PostId: id}

	res, err := c.GetPost(ctx, request)
	if err != nil {
		log.Fatalf("error while calling Post server RPC %v", err)
	}
	s := &models.Post{
		ID:          res.Post.Id,
		Title:       res.Post.Title,
		Description: res.Post.Description,
		Category:    res.Post.Category,
		UserId:      res.Post.UserId,
	}
	log.Printf("response from Post server:%v", res.Post.Title)
	app.render(w, r, "show.page.tmpl", &templateData{
		Post: s,
	})
}
func deletePost(c servicepb1.PostServiceClient) {
	ctx := context.Background()

	request := &servicepb1.DeletePostRequest{PostId: 1}

	response, err := c.DeletePost(ctx, request)
	if err != nil {
		log.Fatalf("error while calling Post server RPC %v", err)
	}
	log.Printf("response from Post server:%v", response.Status)
}
func updatePost(c servicepb1.PostServiceClient) {
	ctx := context.Background()

	request := &servicepb1.UpdatePostRequest{Post: &servicepb1.Post{Id: 1, Title: "changed", Description: "changed_desc", Category: "Tech", UserId: 6}}

	response, err := c.UpdatePost(ctx, request)
	if err != nil {
		log.Fatalf("error while calling Post server RPC %v", err)
	}
	log.Printf("response from Post server:%v", response.Status)
}
func (app *application) getAllPosts(c servicepb1.PostServiceClient, w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := &servicepb1.GetAllPostsRequest{}

	stream, err := c.GetAllPosts(ctx, req)
	if err != nil {
		log.Fatalf("error while calling GET ALL POSTS RPC %v", err)
	}
	defer stream.CloseSend()

	var posts []*models.Post

LOOP:
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break LOOP
			}
			log.Fatalf("error while reciving from get all posts RPC %v", err)
		}
		s := &models.Post{
			ID:          res.Post.Id,
			Title:       res.Post.Title,
			Description: res.Post.Description,
			Category:    res.Post.Category,
			UserId:      res.Post.UserId,
		}
		posts = append(posts, s)
		app.render(w, r, "home.page.tmpl", &templateData{
			Posts: posts,
		})
		log.Printf("response from get all posts:%v \n", res.GetPost().Title)
	}
} //
