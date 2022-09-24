package app

import "github.com/google/uuid"

type User struct {
	UserId        int      `json:"user_id"`
	UserName      string   `json:"user_name"`
	PickupHistory []*Trash `json:"pickup_history"`
	ReportHistory []*Trash `json:"report_history"`
	Rank          int      `json:"rank"`

	// temporary data
	JWTToken      string `json:"jwt_token"`
	FirebaseToken string `json:"firebase_token"`
	// Gamification

	// For reports:
	// each report gets 1 point
	// each confirmation gets 1 point
	// For pickups:
	// The bigger the size, more points
	// The more items, more points
	Score int `json:"score"`
}

// We differentiate hotposts and Trashes by close Coordinates
type Trash struct {
	ID        uuid.UUID `json:"id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	ImageURL  []byte    `json:"image_url"`

	// UserId that reported
	ReportedBy string `json:"reported_by"`
	// Gamification
	ReportNumber int `json:"report_number"`
	Reward       int `json:"reward"`
}
