package app

import "github.com/gorilla/mux"

type Server struct {
	DB     *DB
	Router *mux.Router
}

func (s *Server) RegisterRoutes() {
	s.Router.HandleFunc("/", HelloWorld).Methods("GET")
	s.Router.HandleFunc("/user", GetUser).Methods("GET")
	s.Router.HandleFunc("/user", CreateUser).Methods("POST")
	s.Router.HandleFunc("/trash", ReportTrash).Methods("POST")
	s.Router.HandleFunc("/trash", ConfirmTrash).Methods("SET")
	s.Router.HandleFunc("/trash", PickupTrash).Methods("DELETE")
}
