package ports

import "hexagonal2/core/entity"

type BookingRepository interface {
	GetBookings()([]entity.Booking,error)
	GetABooking(string)(*entity.Booking,error)
	AddBooking(entity.Booking)error
	UpdateBooking(string,entity.Booking)error
	DeleteBooking(string)error
}

type BookingServices interface {
	GetAllBookings()([]entity.BookingRes,error)
	GetBooking(string)(*entity.BookingRes,error)
	AddBooking(entity.Booking)error
	UpdateBooking(string,entity.Booking)error
	DeleteBooking(string)error
}