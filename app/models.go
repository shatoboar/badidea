package app

type User struct {
	UserId        string   `json:"user_id"`
	UserName      string   `json:"user_name"`
	PickupHistory []*Trash `json:"pickup_history"`
	ReportHistory []*Trash `json:"report_history"`

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
	XP   int    `json:"score"`
	Rank string `json:"rank"`
}

// We differentiate hotposts and Trashes by close Coordinates
type Trash struct {
	ID        string  `json:"id"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Image     []byte  `json:"image"`

	// UserId that reported
	ReportedBy string
	// Gamification
	ReportNumber int `json:"report_number"`
	Reward       int `json:"reward"`
}
