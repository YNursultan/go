package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/post/:id", app.getPostAdapter)
	mux.HandleFunc("/post/create", app.createPostAdapter)
	mux.HandleFunc("/post/update", app.updatePostAdapter)
	mux.HandleFunc("/post/delete", app.deletePostAdapter)

	mux.HandleFunc("/user", app.getUserAdapter)
	mux.HandleFunc("/user/:id", app.getUserAdapter)
	mux.HandleFunc("/user/create", app.createUserAdapter)
	mux.HandleFunc("/user/update", app.updateUserAdapter)
	mux.HandleFunc("/user/delete", app.deleteUserAdapter)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
