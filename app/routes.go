package app

import "github.com/gorilla/mux"

type Server struct {
	DB     *DB
	Router *mux.Router
}

func (s *Server) RegisterRoutes() {
	s.Router.HandleFunc("/", HelloWorld).Methods("GET")
	s.Router.HandleFunc("/user", s.GetUser).Methods("GET")
	s.Router.HandleFunc("/user", s.CreateUser).Methods("POST")
	s.Router.HandleFunc("/trash", s.ReportTrash).Methods("POST")
	s.Router.HandleFunc("/trash", s.UpvoteTrash).Methods("SET")
	s.Router.HandleFunc("/trash", s.PickupTrash).Methods("DELETE")
}
