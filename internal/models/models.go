package models

import "time"

// Reservation holds reservation data
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

// Users is the users model
type Users struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// Rooms is the rooms model
type Rooms struct {
	ID        int
	RoomName  string
	UpdatedAt time.Time
	CreatedAt time.Time
}

// Restrictions is the restrictions model
type Restrictions struct {
	ID              int
	RestrictionName string
	UpdatedAt       time.Time
	CreatedAt       time.Time
}

// Reservations is the reservations model
type Reservations struct {
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
	Room      Rooms
}

// RoomRestrictions is the room_restrictions model
type RoomRestrictions struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	UpdatedAt     time.Time
	CreatedAt     time.Time
	Room          Rooms
	Rservation    Reservations
	Restriction   Restrictions
}
