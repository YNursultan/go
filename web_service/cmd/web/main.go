package main

import (
	servicepb1 "final_project/posts_service/api"
	servicepb2 "final_project/user_service/api"
	"flag"
	"github.com/golangcollege/sessions"
	"google.golang.org/grpc"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	c1            servicepb1.PostServiceClient
	c2            servicepb2.UserServiceClient
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	templateCache map[string]*template.Template
}

func main() {
	conn1, err1 := grpc.Dial("localhost:50051", grpc.WithInsecure())
	conn2, err2 := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err1 != nil {
		log.Fatalf("could not connect: %v", err1)
	}
	defer conn1.Close()
	if err2 != nil {
		log.Fatalf("could not connect: %v", err2)
	}
	defer conn2.Close()

	c1 := servicepb1.NewPostServiceClient(conn1)
	c2 := servicepb2.NewUserServiceClient(conn2)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	addr := flag.String("addr", ":4000", "HTTP network address")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		c1:            c1,
		c2:            c2,
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Server running on port %v", *addr)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}
