package repository

import (
	"time"

	"github.com/NomanSalhab/golang_b_n_b_training_project/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(r models.RoomRestriction) error

	SearchAvailabilityByDatesByRoomID(roomID int, start, end time.Time) (bool, error)

	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)

	GetRoomByID(id int) (models.Room, error)
}
