package app

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type DB struct {
	Users       map[uuid.UUID]*User
	Trash       map[uuid.UUID]*Trash
	LeaderBoard map[uuid.UUID]*User
}

func NewDB() *DB {
	return &DB{
		Users:       make(map[uuid.UUID]*User, 0),
		Trash:       make(map[uuid.UUID]*Trash, 0),
		LeaderBoard: make(map[uuid.UUID]*User, 0),
	}
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func getDistance(lat1, lon1, lat2, lon2 float64) float64 {
	R := 6378.137
	dLat := lat2*math.Pi/180 - lat1*math.Pi/180
	dLon := lon2*math.Pi/180 - lon1*math.Pi/180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c

	return d * 1000
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Errorf("Couldn't decode User: %v", err)
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
		log.Errorf("Couldn't decode User: %v", err)
	}

	requestedUser, ok := s.DB.Users[newUser.UserId]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requestedUser)
}

// ReportTrash is called when a User wants to report a new trash.
// If there are trashes in vicinity, we send back the cloesst trashes.
// Otherwise create a new trash
func (s *Server) ReportTrash(w http.ResponseWriter, r *http.Request) {
	var reportedTrash Trash
	err := json.NewDecoder(r.Body).Decode(&reportedTrash)
	if err != nil {
		log.Errorf("Couldn't decode trash: %v", reportedTrash)
		w.WriteHeader(http.StatusBadRequest)
	}
	log.Infof("Got new Trash: %v", reportedTrash)

	uid, err := uuid.NewUUID()
	if err != nil {
		log.Errorf("Failed to generate new uuid: %v", err)
	}

	closestTrashes := make([]*Trash, 0)
	for _, trash := range s.DB.Trash {
		distanceInMeter := getDistance(trash.Latitude, trash.Longitude, reportedTrash.Latitude, reportedTrash.Longitude)
		if distanceInMeter < 15 {
			closestTrashes = append(closestTrashes, trash)
		}
	}

	// if there are no close trash, then we can add this as a new trash
	if len(closestTrashes) == 0 {
		reportedTrash.ID = uid
		reportedTrash.ReportNumber = 1
		reportedTrash.Reward = 1
		s.DB.Trash[uid] = &reportedTrash
		w.WriteHeader(http.StatusOK)
		return
	}

	json.NewEncoder(w).Encode(closestTrashes)
}

// Confirms Trash exists. User gets a point for the upvote
func (s *Server) UpvoteTrash(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func (s *Server) CreateNewTrash(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) PickupTrash(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}
