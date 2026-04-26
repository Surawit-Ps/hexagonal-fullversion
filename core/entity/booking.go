package entity

import (
	"time"
)

type Booking struct {
    ID        string
    BookingID string
    RoomID     string
    CustomerID   string
    Service   string
    StartTime time.Time
    EndTime   time.Time
    Status    string // pending, confirmed, cancelled
    CreatedAt time.Time
}

type BookingRes struct {
    BookingID string
	RoomID     string
	CustomerID   string
    Service   string
	StartTime time.Time
	EndTime   time.Time
	Status    string // pending, confirmed, cancelled
}