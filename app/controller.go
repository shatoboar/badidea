package app

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	ReportReward = 1
	BaseReward   = 5
)

type DB struct {
	Users       map[string]*User
	Trash       map[uuid.UUID]*Trash
	LeaderBoard map[uuid.UUID]*User
}

func initMockTrash() map[uuid.UUID]*Trash {
	trash := make(map[uuid.UUID]*Trash, 0)

	uid, _ := uuid.NewUUID()
	trash[uid] = &Trash{
		ID:           uid,
		Latitude:     52.494121,
		Longitude:    13.445063,
		ImageURL:     "https://www.stadtbetrieb-frechen.de/storage/media/images/209/conversions/sperrmuell-slide.jpg",
		ReportedBy:   "gilles",
		ReportNumber: 1,
		Reward:       4,
	}

	uid, _ = uuid.NewUUID()
	trash[uid] = &Trash{
		ID:           uid,
		Latitude:     52.490249,
		Longitude:    13.437251,
		ImageURL:     "https://umziehen.de/media/cache/article_image/cms/2018/12/Sperrmuell-entsorgen-Umziehen-coramueller-iStock.jpg?869457",
		ReportedBy:   "karsten",
		ReportNumber: 5,
		Reward:       10,
	}

	uid, _ = uuid.NewUUID()
	trash[uid] = &Trash{
		ID:           uid,
		Latitude:     52.497118,
		Longitude:    13.434719,
		ImageURL:     "https://www.zvo.com/files/images/3-entsorgung/sperrmuellabholung/sperrmuell-bereitgestellt.jpg",
		ReportedBy:   "filip",
		ReportNumber: 3,
		Reward:       7,
	}

	uid, _ = uuid.NewUUID()
	trash[uid] = &Trash{
		ID:           uid,
		Latitude:     52.492664,
		Longitude:    13.461477,
		ImageURL:     "https://www.ruempelmannschaft.de/wp-content/uploads/2022/06/sperrmuell-abholung-koeln.jpg",
		ReportedBy:   "mantas",
		ReportNumber: 1,
		Reward:       10,
	}

	uid, _ = uuid.NewUUID()
	trash[uid] = &Trash{
		ID:           uid,
		Latitude:     52.492939,
		Longitude:    13.452390,
		ImageURL:     "https://www.avea.info/images/titel/fotolia_110482889_l_sperrmuell_151_md.jpg",
		ReportedBy:   "daniel",
		ReportNumber: 8,
		Reward:       10,
	}
	return trash

}

func initMockUsers() map[string]*User {
	users := make(map[string]*User, 0)
	users["gilles"] = &User{
		UserName:      "gilles",
		PickupHistory: []*Trash{},
		ReportHistory: []*Trash{},
	}
	users["daniel"] = &User{
		UserName:      "daniel",
		PickupHistory: []*Trash{},
		ReportHistory: []*Trash{},
	}
	users["mantas"] = &User{
		UserName:      "mantas",
		PickupHistory: []*Trash{},
		ReportHistory: []*Trash{},
	}
	users["filip"] = &User{
		UserName:      "filip",
		PickupHistory: []*Trash{},
		ReportHistory: []*Trash{},
	}
	users["karsten"] = &User{
		UserName:      "karsten",
		PickupHistory: []*Trash{},
		ReportHistory: []*Trash{},
	}
	return users
}

func NewDB() *DB {
	return &DB{
		Users:       initMockUsers(),
		Trash:       initMockTrash(),
		LeaderBoard: make(map[uuid.UUID]*User, 0),
	}
}

func enableCors(w *http.ResponseWriter) {
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func getDistance(lat1, lon1, lat2, lon2 float64) float64 {
	R := 6378.137
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c

	return d * 1000
}

// Lookup firebase token to check whether this is valid
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Errorf("Couldn't decode User: %v", err)
	}

	_, ok := s.DB.Users[newUser.UserName]
	if ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.DB.Users[newUser.UserName] = &newUser
	log.Infof("A new user was added to the DB %v", newUser)
	w.WriteHeader(http.StatusCreated)
}

// Requesting detailed Userdata
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	userName, err := decodeUserName(r)
	if err != nil {
		log.Errorf("Failed to decode userName: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestedUser, ok := s.DB.Users[userName]
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
	enableCors(&w)
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
		if int(distanceInMeter) < 15 {
			closestTrashes = append(closestTrashes, trash)
		}
	}

	// there are other options, give user opportunity to decide
	if len(closestTrashes) > 0 {
		json.NewEncoder(w).Encode(closestTrashes)
		return
	}

	userName, err := decodeUserName(r)
	if err != nil {
		log.Errorf("Failed to get userName: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, existing := s.DB.Users[userName]
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
	reportedTrash.Reward = 5
	reportedTrash.ReportedBy = user.UserName
	s.DB.Trash[uid] = &reportedTrash
	w.WriteHeader(http.StatusCreated)

}

// Confirms Trash exists. User gets a point for the upvote
func (s *Server) UpvoteTrash(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
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

	userName, err := decodeUserName(r)
	if err != nil {
		log.Errorf("Failed to get userName: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	trash.ReportNumber++
	trash.Reward += ReportReward + BaseReward

	user, ok := s.DB.Users[userName]
	if !ok {
		log.Errorf("User doesn't exist: %v", err)
	}
	user.ReportHistory = append(user.ReportHistory, trash)
	user.Score += ReportReward
	w.WriteHeader(http.StatusOK)
}

func (s *Server) CreateNewTrash(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
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

	userName, err := decodeUserName(r)
	if err != nil {
		log.Errorf("Failed to get userName: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newTrash.ID = uid
	newTrash.ReportNumber = 1
	newTrash.Reward = BaseReward
	newTrash.ReportedBy = userName
	s.DB.Trash[uid] = &newTrash

	user, ok := s.DB.Users[userName]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.ReportHistory = append(user.ReportHistory, &newTrash)
	user.Score += ReportReward

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) PickupTrash(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var pickedTrash Trash
	err := json.NewDecoder(r.Body).Decode(&pickedTrash)
	if err != nil {
		log.Errorf("Couldn't decode trash: %v", pickedTrash)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userName, err := decodeUserName(r)
	if err != nil {
		log.Errorf("Failed to get userName: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := s.DB.Users[userName]
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
	if len(s.DB.Trash) < 5 {
		newTrash := s.createMockTrash()
		s.DB.Trash[newTrash.ID] = newTrash

		newUser := s.DB.Users[newTrash.ReportedBy]

		user.ReportHistory = append(newUser.ReportHistory, newTrash)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetTrash(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	allTrashes := make([]*Trash, 0)
	for _, val := range s.DB.Trash {
		allTrashes = append(allTrashes, val)
	}

	json.NewEncoder(w).Encode(allTrashes)
}

func (s *Server) createMockTrash() *Trash {
	names := [5]string{"gilles", "daniel", "mantas", "filip", "karsten"}
	images := [5]string{"https://www.ra-kotz.de/wp-content/uploads/2018/07/bigstock-183494932.jpg",
		"https://fotos.verwaltungsportal.de/news/7/4/7/3/3/4/gross/20220701_Sperrmuell.jpg",
		"https://www.gea.de/cms_media/module_img/80116/40058165_1_detail_PEL_Sperrmuell.jpg",
		"https://www.cannstatter-zeitung.de/media.media.25cbf81f-7997-4cfd-a6d4-e0ad49ea786e.original1024.jpg",
		"https://www.berliner-mieterverein.de/uploads/2017/05/061724-a-altes-sofa.jpg"}
	name := names[rand.Intn(5)]
	image := images[rand.Intn(5)]

	uid, err := uuid.NewUUID()
	if err != nil {
		log.Errorf("Failed to generate new uuid: %v", err)
	}

	return &Trash{
		ID:           uid,
		ReportNumber: 1,
		Reward:       BaseReward,
		ReportedBy:   name,
		Latitude:     52.497116 + float64(rand.Intn(8000)/1000000),
		Longitude:    13.434719 + float64(rand.Intn(8000)/1000000),
		ImageURL:     image,
	}
}

func (s *Server) GetLeaderBoard(w http.ResponseWriter, r *http.Request) {
	leaderBoard := getTopUsers(s.DB.Users, 10)
	length := len(leaderBoard)
	leader := 1
	for i := length - 1; i > 0; i-- {
		user := s.DB.Users[leaderBoard[i].UserName]
		s.DB.Users[user.UserName].Rank = leader
		leader++
		if i == length-9 {
			break
		}
	}

	json.NewEncoder(w).Encode(leaderBoard)
}
