package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	e "hexagonal2/pkg/errors"
	"hexagonal2/core/entity"
)

type subRepositoryMongo struct {
	col *mongo.Collection
}

type SubMongo struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         string `bson:"user_id" json:"user_id"`
	SubscriptionID string `bson:"subscription_id" json:"subscription_id"`
	Status         string `bson:"status" json:"status"`
	ExpiryDate     time.Time `bson:"expiry_date" json:"expiry_date"`
}

func subEnToMongo(s entity.Subscription) SubMongo {
	return SubMongo{
		ID:     primitive.NewObjectID(),
		UserID: s.UserID,
		SubscriptionID: s.SubscriptionID,
		Status: s.Status,
		ExpiryDate: s.ExpiryDate,
	}
}

func subMongoToEn(m SubMongo) entity.Subscription {
	return entity.Subscription{
		ID: m.ID.Hex(),
		UserID: m.UserID,
		SubscriptionID: m.SubscriptionID,
		Status: m.Status,
		ExpiryDate: m.ExpiryDate,
	}
}

func NewSubRepositoryMongo(client *mongo.Client, dbName string) *subRepositoryMongo {
	return &subRepositoryMongo{col: client.Database(dbName).Collection("subscriptions")}
}

func (r *subRepositoryMongo) CreateSubscription(subscription entity.Subscription) error {
	subscription.ExpiryDate = time.Now().AddDate(0, 1, 0)
	subMongo := subEnToMongo(subscription)
	_, err := r.col.InsertOne(context.Background(), subMongo)
	if err != nil {
		return e.ErrInternalServer
	}
	return nil
}

func (r *subRepositoryMongo) GetSubscriptionByUserID(userID string) (*entity.Subscription, error) {
	var sub SubMongo
	err := r.col.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&sub)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, e.ErrNotFound
		}
		return nil, e.ErrInternalServer
	}
	subscription := subMongoToEn(sub)
	return &subscription, nil
}