package app

import "github.com/google/uuid"

type User struct {
	UserId        int      `json:"user_id"`
	UserName      string   `json:"user_name"`
	PickupHistory []*Trash `json:"pickup_history"`
	ReportHistory []*Trash `json:"report_history"`

	// temporary data
	// JWTToken      string `json:"jwt_token"`
	// FirebaseToken string `json:"firebase_token"`
	// Gamification

	// For reports:
	// each report gets 1 point
	// each confirmation gets 1 point
	// For pickups:
	// The bigger the size, more points
	// The more items, more points
	//Rank Rank `json:"rank"`
	Rank  int
	Score int
	Title string
}

// We differentiate hotposts and Trashes by close Coordinates
type Trash struct {
	ID        uuid.UUID `json:"id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	ImageURL  string    `json:"image_url"`

	// UserId that reported
	ReportedBy int `json:"reported_by"`
	// Gamification
	ReportNumber int `json:"report_number"`
	Reward       int `json:"reward"`
}
