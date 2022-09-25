package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/slash/badidea/app"
)

const port = ":8080"

func main() {
	// client := initFirebase()
	// auth, err := client.Auth(context.Background())
	// if err != nil {
	// 	fmt.Printf("Failed to create a new service account: %v", err)
	// }
	router := mux.NewRouter()
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE"})
	// ttl := handlers.MaxAge(3600)
	origins := handlers.AllowedOrigins([]string{"*"})

	server := &app.Server{
		DB:     app.NewDB(),
		Router: router,
		// Auth: &app.AuthClient{
		// 	Client: auth,
		// },
	}

	http.Handle("/", server.Router)
	server.RegisterRoutes()
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, handlers.CORS(credentials, methods, origins)(router)))
}

// func initFirebase() *firebase.App {
// 	opt := option.WithCredentialsFile("./serviceAccount.json")
// 	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
// 	if err != nil {
// 		fmt.Printf("Failed to initialize firebase auth: %v", err)
// 		os.Exit(1)
// 	}
// 	return firebaseApp
// }
