package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	DB     *DB
	Router *mux.Router
	// Auth   Verifier
}

type Verifier interface {
	verifyUser(r *http.Request) bool
}

func (s *Server) RegisterRoutes() {
	s.Router.HandleFunc("/", HelloWorld).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/user", s.GetUser).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/user", s.CreateUser).Methods("POST", "OPTIONS")
	s.Router.HandleFunc("/trash", s.GetTrash).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/trash", s.ReportTrash).Methods("POST", "OPTIONS")
	s.Router.HandleFunc("/trash/new", s.CreateNewTrash).Methods("POST", "OPTIONS")
	s.Router.HandleFunc("/trash", s.UpvoteTrash).Methods("PUT", "OPTIONS")
	s.Router.HandleFunc("/trash", s.PickupTrash).Methods("DELETE", "OPTIONS")
}
