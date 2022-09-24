package app

import (
	"fmt"
	"net/http"
)

type DB struct {
	Users       []*User
	Trash       []*Trash
	LeaderBoard []*User
}

func NewDB() *DB {
	return &DB{
		Users:       make([]*User, 0),
		Trash:       make([]*Trash, 0),
		LeaderBoard: make([]*User, 0),
	}
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Got user")
}

// Report Trash creates a new Trash.
func ReportTrash(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Heyllo world")
}

// Confirms Trash exists. User gets a point for the upvote
func ConfirmTrash(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func PickupTrash(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}
