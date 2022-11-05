package models

import "time"

// // Reservation holds reservation data
// type Reservation struct {
// 	FirstName string
// 	LastName  string
// 	Email     string
// 	Phone     string
// }

// User is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// Room is the room model
type Room struct {
	ID        int
	RoomName  string
	UpdatedAt time.Time
	CreatedAt time.Time
}

// Restriction is the restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	UpdatedAt       time.Time
	CreatedAt       time.Time
}

// Reservation is the reservation model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int //* (Or) Room Rooms to inckude all room information
	UpdatedAt time.Time
	CreatedAt time.Time
	Room      Room
}

// RoomRestriction is the room_restriction model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	UpdatedAt     time.Time
	CreatedAt     time.Time
	Room          Room
	Rservation    Reservation
	Restriction   Restriction
}
