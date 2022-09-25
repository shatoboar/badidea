package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
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
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                            // All origins
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"}, // Allowing only get, just an example
		AllowedHeaders: []string{"*"},                            // Allowing only get, just an example
	})

	server := &app.Server{
		DB:     app.NewDB(),
		Router: router,
		// Auth: &app.AuthClient{
		// 	Client: auth,
		// },
	}

	server.RegisterRoutes()
	http.Handle("/", server.Router)
	log.Infof("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, c.Handler(router)))
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
