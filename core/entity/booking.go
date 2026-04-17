package entity

import (
	"time"
)

type Booking struct {
    ID        string
    BookingID string
    PetID     string
    OwnerID   string
    Service   string
    StartTime time.Time
    EndTime   time.Time
    Status    string // pending, confirmed, cancelled
    CreatedAt time.Time
}

type BookingRes struct {
    BookingID string
	PetID     string
	OwnerID   string
    Service   string
	StartTime time.Time
	EndTime   time.Time
	Status    string // pending, confirmed, cancelled
}