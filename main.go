package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"github.com/slash/badidea/app"
	"google.golang.org/api/option"
)

const port = ":8080"

func main() {
	client := initFirebase()
	auth, err := client.Auth(context.Background())
	if err != nil {
		fmt.Printf("Failed to create a new service account: %v", err)
	}

	server := &app.Server{
		DB:     app.NewDB(),
		Router: mux.NewRouter(),
		Auth: &app.AuthClient{
			Client: auth,
		},
	}

	http.Handle("/", server.Router)
	server.RegisterRoutes()
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func initFirebase() *firebase.App {
	opt := option.WithCredentialsFile("./serviceAccount.json")
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("Failed to initialize firebase auth: %v", err)
		os.Exit(1)
	}
	return firebaseApp
}
