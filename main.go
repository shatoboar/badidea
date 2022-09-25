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
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	router := mux.NewRouter()

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
	log.Fatal(http.ListenAndServe(port, handlers.CORS(headersOk, originsOk, methodsOk)(router)))
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
