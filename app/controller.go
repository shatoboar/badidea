package app

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const ReportReward = 1

type DB struct {
	Users       map[int]*User
	Trash       map[uuid.UUID]*Trash
	LeaderBoard map[uuid.UUID]*User
}

func NewDB() *DB {
	return &DB{
		Users:       make(map[int]*User, 0),
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

// Lookup firebase token to check whether this is valid
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Errorf("Couldn't decode User: %v", err)
	}

	_, ok := s.DB.Users[newUser.UserId]
	if ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.DB.Users[newUser.UserId] = &newUser
	log.Infof("A new user was added to the DB %v", newUser)
	w.WriteHeader(http.StatusCreated)
}

// Requesting detailed Userdata
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	if !verifyUser(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := decodeUserID(r)
	if err != nil {
		log.Errorf("Failed to decode userID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestedUser, ok := s.DB.Users[userID]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requestedUser)
}

// ReportTrash is called when a User wants to report a new trash.
// If there are trashes in vicinity, we send back the closest trashes.
// Otherwise create a new trash
func (s *Server) ReportTrash(w http.ResponseWriter, r *http.Request) {
	if !verifyUser(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var reportedTrash Trash
	err := json.NewDecoder(r.Body).Decode(&reportedTrash)
	if err != nil {
		log.Errorf("Couldn't decode trash: %v", reportedTrash)
		w.WriteHeader(http.StatusBadRequest)
	}
	log.Infof("Got new Trash: %v", reportedTrash)

	closestTrashes := make([]*Trash, 0)
	for _, trash := range s.DB.Trash {
		distanceInMeter := getDistance(trash.Latitude, trash.Longitude, reportedTrash.Latitude, reportedTrash.Longitude)
		if distanceInMeter < 15 {
			closestTrashes = append(closestTrashes, trash)
		}
	}

	// there are other options, give user opportunity to decide
	if len(closestTrashes) > 0 {
		json.NewEncoder(w).Encode(closestTrashes)
		return
	}

	userID, err := decodeUserID(r)
	if err != nil {
		log.Errorf("Failed to get userID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, existing := s.DB.Users[userID]
	if !existing {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.ReportHistory = append(user.ReportHistory, &reportedTrash)
	user.Score += ReportReward
	// TODO: user.Rank

	uid, err := uuid.NewUUID()
	if err != nil {
		log.Errorf("Failed to generate new uuid: %v", err)
	}
	reportedTrash.ID = uid
	reportedTrash.ReportNumber = 1
	reportedTrash.Reward = 1
	s.DB.Trash[uid] = &reportedTrash
	w.WriteHeader(http.StatusCreated)

}

// Confirms Trash exists. User gets a point for the upvote
func (s *Server) UpvoteTrash(w http.ResponseWriter, r *http.Request) {
	if !verifyUser(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var existingTrash Trash
	err := json.NewDecoder(r.Body).Decode(&existingTrash)
	if err != nil {
		log.Errorf("Couldn't decode trash: %v", existingTrash)
		w.WriteHeader(http.StatusBadRequest)
	}
	log.Infof("Got new Trash: %v", existingTrash)

	trash, existing := s.DB.Trash[existingTrash.ID]
	if !existing {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := decodeUserID(r)
	if err != nil {
		log.Errorf("Failed to get userID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	trash.ReportNumber++

	user, ok := s.DB.Users[userID]
	if !ok {
		log.Errorf("User doesn't exist: %v", err)
	}
	user.ReportHistory = append(user.ReportHistory, trash)
	user.Score += ReportReward
	w.WriteHeader(http.StatusOK)
}

func (s *Server) CreateNewTrash(w http.ResponseWriter, r *http.Request) {
	if !verifyUser(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var newTrash Trash
	err := json.NewDecoder(r.Body).Decode(&newTrash)
	if err != nil {
		log.Errorf("Couldn't decode trash: %v", newTrash)
		w.WriteHeader(http.StatusBadRequest)
	}
	log.Infof("Got new Trash: %v", newTrash)

	uid, err := uuid.NewUUID()
	if err != nil {
		log.Errorf("Failed to generate new uuid: %v", err)
	}

	userID, err := decodeUserID(r)
	if err != nil {
		log.Errorf("Failed to get userID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newTrash.ID = uid
	newTrash.ReportNumber = 1
	newTrash.Reward = 1
	newTrash.ReportedBy = userID
	s.DB.Trash[uid] = &newTrash

	user, ok := s.DB.Users[userID]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.ReportHistory = append(user.ReportHistory, &newTrash)

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) PickupTrash(w http.ResponseWriter, r *http.Request) {
	if !verifyUser(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var pickedTrash Trash
	err := json.NewDecoder(r.Body).Decode(&pickedTrash)
	if err != nil {
		log.Errorf("Couldn't decode trash: %v", pickedTrash)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := decodeUserID(r)
	if err != nil {
		log.Errorf("Failed to get userID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := s.DB.Users[userID]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.Score += pickedTrash.Reward
	user.PickupHistory = append(user.PickupHistory, &pickedTrash)

	log.Infof("Decoded trash: %v", pickedTrash)

	_, ok = s.DB.Trash[pickedTrash.ID]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// verify that we can pick up

	delete(s.DB.Trash, pickedTrash.ID)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetTrash(w http.ResponseWriter, r *http.Request) {
	if !verifyUser(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	allTrashes := make([]*Trash, 0)
	for _, val := range s.DB.Trash {
		allTrashes = append(allTrashes, val)
	}

	json.NewEncoder(w).Encode(allTrashes)
}
