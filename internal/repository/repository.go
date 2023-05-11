package repository

import "github.com/MarynaMarkova/Go-bookings/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) error
}