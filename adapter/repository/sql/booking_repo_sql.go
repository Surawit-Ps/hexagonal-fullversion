package sql

import ("hexagonal2/core/entity"
"gorm.io/gorm"
"github.com/google/uuid"
"fmt"
e "hexagonal2/pkg/errors"
"time"
"strings"
"strconv")

type bookingRepositoryDB struct{
	db *gorm.DB
}

func NewBookingRepositoryDB(db *gorm.DB) bookingRepositoryDB{
	return bookingRepositoryDB{db:db}
}

type Booking struct{
	ID string `gorm:"primaryKey"`
	BookingID string
	RoomID  string
	CustomerID string
	Service string
	StartDate time.Time
	EndDate time.Time
	CreatedAt time.Time
}

func EnToGormBooking(b entity.Booking)Booking{
	return Booking{
		ID: b.ID,
		BookingID: b.BookingID,
		RoomID: b.RoomID,
		CustomerID: b.CustomerID,
		Service: b.Service,
		StartDate: b.StartTime,
		EndDate: b.EndTime,
		CreatedAt: b.CreatedAt,
	}
}

func GormToEnBooking(b Booking)entity.Booking{
	return entity.Booking{
		ID: b.ID,
		BookingID: b.BookingID,
		RoomID: b.RoomID,
		CustomerID: b.CustomerID,
		Service: b.Service,
		StartTime: b.StartDate,
		EndTime: b.EndDate,
		CreatedAt: b.CreatedAt,
	}
}

func (r bookingRepositoryDB)GetBookings()([]entity.Booking,error){
	var bookings []Booking
	result := r.db.Find(&bookings)
	if result.Error != nil {
		return nil, e.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return nil, e.ErrNotFound
	}
	var bookingEn []entity.Booking
	for _, b := range bookings {
		bookingEn = append(bookingEn, GormToEnBooking(b))
	}
	return bookingEn,nil
}

func (r bookingRepositoryDB)GetABooking(id string)(*entity.Booking,error){
	var booking Booking
	result := r.db.Where("id = ?", id).First(&booking)
	if result.Error != nil {
		return nil, e.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return nil, e.ErrNotFound
	}
	bookingEn := GormToEnBooking(booking)
	return &bookingEn,nil
}

func (r bookingRepositoryDB) GenerateID() string {
	var bookingID entity.Booking
	today := time.Now().Format("20060102")
	seq := 1
	err := r.db.Last(&bookingID)
	if err == nil {
		num := strings.Split(bookingID.BookingID, "-")
		if len(num) == 2 {
			if n, parseErr := strconv.Atoi(num[1]); parseErr == nil {
				seq = n + 1
			}
		}
	}
	return fmt.Sprintf("PB%s-%03d", today, seq)
}

func (r bookingRepositoryDB)AddBooking(b entity.Booking)error{
	b.ID = uuid.New().String()
	b.BookingID = r.GenerateID()
	b.CreatedAt = time.Now()
	bookingModel := EnToGormBooking(b)
	result := r.db.Create(&bookingModel)
	if result.Error != nil {
		return e.ErrInternalServer
	}
	return nil
}

func (r bookingRepositoryDB)UpdateBooking(id string,b entity.Booking)error{
	var booking Booking
	result := r.db.Where("id = ?", id).First(&booking)
	if result.Error != nil {
		return e.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return e.ErrNotFound
	}	
	booking.StartDate = b.StartTime
	booking.EndDate = b.EndTime
	result = r.db.Save(&booking)
	if result.Error != nil {
		return e.ErrInternalServer
	}
	return nil
}

func (r bookingRepositoryDB)DeleteBooking(id string)error{
	result := r.db.Delete(&Booking{},"id = ?",id)	
	if result.Error != nil {
		return e.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return e.ErrNotFound
	}
	return nil
}