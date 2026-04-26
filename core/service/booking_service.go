package service

import (
"hexagonal2/core/ports"
"hexagonal2/core/entity")

type BookingService struct {
	bookingRepo ports.BookingRepository
}

func NewBookingService(bookingRepo ports.BookingRepository) *BookingService {
	return &BookingService{bookingRepo: bookingRepo}
}

func (s *BookingService) GetAllBookings() ([]entity.BookingRes, error) {
	bookings, err := s.bookingRepo.GetBookings()
	if err != nil {
		return nil, err
	}
	var bookingRes []entity.BookingRes
	for _, b := range bookings {
		br := entity.BookingRes{
			BookingID: b.BookingID,
			CustomerID: b.CustomerID,
			RoomID:     b.RoomID,
			Service:    b.Service,
			StartTime:  b.StartTime,
			EndTime:    b.EndTime,
			Status:    b.Status,

		}
		bookingRes = append(bookingRes, br)
	}
	return bookingRes, nil
}

func (s *BookingService) GetBooking(id string) (*entity.BookingRes, error) {
	booking, err := s.bookingRepo.GetABooking(id)
	if err != nil {
		return nil, err
	}
	br := entity.BookingRes{
		BookingID: booking.BookingID,
		CustomerID: booking.CustomerID,
		RoomID:     booking.RoomID,
		StartTime:  booking.StartTime,
		EndTime:    booking.EndTime,
		Status:    booking.Status,
	}
	return &br, nil
}

func (s *BookingService) AddBooking(booking entity.Booking) error {
	return s.bookingRepo.AddBooking(booking)
}	

func (s *BookingService) UpdateBooking(id string, booking entity.Booking) error {
	return s.bookingRepo.UpdateBooking(id, booking)
}

func (s *BookingService) DeleteBooking(id string) error {
	return s.bookingRepo.DeleteBooking(id)
}

