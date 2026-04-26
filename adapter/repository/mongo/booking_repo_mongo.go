package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hexagonal2/core/entity"
	e "hexagonal2/pkg/errors"
	"strconv"
	"fmt"
)



type bookingRepositoryMongo struct {
	col *mongo.Collection
}

type BookingMongo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookingID  string `bson:"booking_id" json:"booking_id"`
	RoomID     string `bson:"room_id" json:"room_id"`
	CustomerID     string `bson:"customer_id" json:"customer_id"`
	StartDate time.Time `bson:"start_date" json:"start_date"`
	Service    string `bson:"service" json:"service"`
	EndDate   time.Time `bson:"end_date" json:"end_date"`
	Status    string `bson:"status" json:"status"` // Booked, Completed, Cancelled
    CreatedAt  time.Time
}

func bookingEnToMongo(b entity.Booking) BookingMongo {
	return BookingMongo{
		ID:        primitive.NewObjectID(),
		CustomerID:     b.CustomerID,
		RoomID:     b.RoomID,
		StartDate: b.StartTime,
		EndDate:   b.EndTime,
		Status:    b.Status,
		Service:   b.Service,
		CreatedAt: time.Now(),
	}
}

func bookingMongoToEn(m BookingMongo) entity.Booking {
	return entity.Booking{
		ID:        m.ID.Hex(),
		BookingID: m.BookingID,
		CustomerID: m.CustomerID,
		RoomID:     m.RoomID,
		StartTime: m.StartDate,
		EndTime:   m.EndDate,
		Status:    m.Status,
		Service:   m.Service,
		CreatedAt: m.CreatedAt,
	}
}

func NewBookingRepositoryMongo(client *mongo.Client, dbName string) *bookingRepositoryMongo {
	return &bookingRepositoryMongo{col: client.Database(dbName).Collection("bookings")}
}

func (r *bookingRepositoryMongo) GetBookings() ([]entity.Booking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()	
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, e.ErrInternalServer
	}
	defer cur.Close(ctx)

	var out []entity.Booking
	for cur.Next(ctx) {
		var b BookingMongo
		err := cur.Decode(&b)
		if err != nil {
			return nil, e.ErrInternalServer
		}
		out = append(out, bookingMongoToEn(b))
	}	
	return out, nil
}

func (r *bookingRepositoryMongo) GetABooking(id string) (*entity.Booking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()	
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, e.ErrInvalidID
	}
	var b BookingMongo
	err = r.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&b)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, e.ErrNotFound
		}
		return nil, e.ErrInternalServer
	}
	en := bookingMongoToEn(b)
	return &en, nil
}

func (r *bookingRepositoryMongo) AddBooking(b entity.Booking) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	b.ID = primitive.NewObjectID().Hex() // Generate new ID for MongoDB
	b.BookingID = r.GenerateBookingID() // Generate booking ID
	bookingMongo := bookingEnToMongo(b)
	_, err := r.col.InsertOne(ctx, bookingMongo)
	if err != nil {
		return e.ErrInternalServer
	}
	return nil
}

func (r *bookingRepositoryMongo) UpdateBookingStatus(id string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()	
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return e.ErrInvalidID
	}
	update := bson.M{"$set": bson.M{"status": status}}
	result, err := r.col.UpdateOne(ctx,
		bson.M{"_id": objID},
		update)
	if err != nil {
		return e.ErrInternalServer
	}
	if result.MatchedCount == 0 {
		return e.ErrNotFound
	}
	return nil
}

func (r *bookingRepositoryMongo) UpdateBooking(id string, b entity.Booking) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()	
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return e.ErrInvalidID
	}
	update := bson.M{
		"$set": bson.M{
			"booking_id": b.BookingID,
			"customer_id": b.CustomerID,
			"room_id": b.RoomID,
			"start_date": b.StartTime,
			"end_date": b.EndTime,
			"status":     b.Status,
			"service":    b.Service,
			"created_at": b.CreatedAt,
		},
	}
	result, err := r.col.UpdateOne(ctx,
		bson.M{"_id": objID},
		update)
	if err != nil {
		return e.ErrInternalServer
	}
	if result.MatchedCount == 0 {
		return e.ErrNotFound
	}
	return nil
}

func (r *bookingRepositoryMongo) DeleteBooking(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()	
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return e.ErrInvalidID
	}
	result, err := r.col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return e.ErrInternalServer
	}
	if result.DeletedCount == 0 {
		return e.ErrNotFound
	}
	return nil
}

func (r *bookingRepositoryMongo) GenerateBookingID() string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var booking entity.Booking
	opts := options.FindOne().SetSort(bson.D{{"_id", -1}})
	err := r.col.FindOne(ctx, bson.M{}, opts).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "BKG-000001"
		}
		return "BKG-000001" // fallback in case of error
	}
	today := time.Now().Format("20060102")
	lastID := booking.BookingID
	lastSeqStr := lastID[12:] // Extract sequence part
	lastSeq, err := strconv.Atoi(lastSeqStr)
	if err != nil {
		lastSeq = 0 // default to 0 if parsing fails
	}
	newSeq := lastSeq + 1
	return fmt.Sprintf("BKG-%s-%06d", today, newSeq)
}



