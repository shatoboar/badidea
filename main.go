package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/slash/badidea/app"
)

const port = ":8080"

func main() {
	server := &app.Server{
		DB:     app.NewDB(),
		Router: mux.NewRouter(),
	}

	http.Handle("/", server.Router)
	server.RegisterRoutes()
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
