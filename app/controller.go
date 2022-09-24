package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type DB struct {
	Users       map[string]*User
	Trash       map[string]*Trash
	LeaderBoard map[string]*User
}

func NewDB() *DB {
	return &DB{
		Users:       make(map[string]*User, 0),
		Trash:       make(map[string]*Trash, 0),
		LeaderBoard: make(map[string]*User, 0),
	}
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Fatalf("Couldn't decode User: %v", err)
	}

	_, ok := s.DB.Users[newUser.UserId]
	if !ok {
		s.DB.Users[newUser.UserId] = &newUser
	}

	log.Infof("A new user was added to the DB %v", newUser)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Fatalf("Couldn't decode User: %v", err)
	}

	requestedUser, ok := s.DB.Users[newUser.UserId]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requestedUser)
}

// Report Trash creates a new Trash.
func (s *Server) ReportTrash(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Heyllo world")
}

// Confirms Trash exists. User gets a point for the upvote
func (s *Server) ConfirmTrash(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func (s *Server) PickupTrash(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}
