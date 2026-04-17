package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	e "hexagonal2/pkg/errors"
	"hexagonal2/core/entity"
)

type dogsRepositoryMongo struct {
	col *mongo.Collection
}

type DogMongo struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string `bson:"name" json:"name"`
	Age     uint   `bson:"age" json:"age"`
	Colour  string `bson:"colour" json:"colour"`
	UserID string `bson:"user_id" json:"user_id"`
}

func dogEnToMongo(d entity.Dogs) DogMongo {
	return DogMongo{
		ID:      primitive.NewObjectID(),
		Name:    d.Name,
		Age:     d.Age,
		Colour:  d.Colour,
		UserID: d.UserID,
	}
}

func dogMongoToEn(m DogMongo) entity.Dogs {
	return entity.Dogs{
		Id:      m.ID.Hex(),		
		Name:    m.Name,
		Age:     m.Age,
		Colour:  m.Colour,
		UserID: m.UserID,
	}
}

func NewDogsRepositoryMongo(client *mongo.Client, dbName string) *dogsRepositoryMongo {
	return &dogsRepositoryMongo{col: client.Database(dbName).Collection("dogs")}
}

func (r *dogsRepositoryMongo) GetDogs() ([]entity.Dogs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, e.ErrInternalServer
	}
	defer cur.Close(ctx)

	var out []entity.Dogs
	for cur.Next(ctx) {
		var m DogMongo
		if err := cur.Decode(&m); err != nil {
			return nil, e.ErrInternalServer
		}
		out = append(out, dogMongoToEn(m))
	}
	return out, nil
}

func (r *dogsRepositoryMongo) GetADogs(id string) (*entity.Dogs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var m DogMongo
	if err := r.col.FindOne(ctx, bson.M{"id": id}).Decode(&m); err != nil {
		return nil, e.ErrDogNotFound
	}
	en := dogMongoToEn(m)
	return &en, nil
}

func (r *dogsRepositoryMongo) AddDog(d entity.Dogs, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if d.Id == "" {
		d.Id = primitive.NewObjectID().Hex()
	}
	d.UserID = userID

	m := dogEnToMongo(d)
	filter := bson.M{"id": m.ID}
	update := bson.M{"$setOnInsert": m}
	opts := options.Update().SetUpsert(true)
	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	return err
}
