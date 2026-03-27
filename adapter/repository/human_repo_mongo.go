package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"hgo/core/entity"
)

type humanRepositoryMongo struct {
	col *mongo.Collection
}

type HumanMongo struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string `bson:"name" json:"name"`
	LastName string `bson:"last_name" json:"last_name"`
	Age      int    `bson:"age" json:"age"`
	Email    string `bson:"email" json:"email"`
	Tel      string `bson:"tel" json:"tel"`
}

func humanEnToMongo(h entity.Humans) HumanMongo {
	return HumanMongo{
		ID:       primitive.NewObjectID(),
		Name:     h.Name,
		LastName: h.LastName,
		Age:      h.Age,
		Email:    h.Email,
		Tel:      h.Tel,
	}
}

func humanMongoToEn(m HumanMongo) entity.Humans {
	return entity.Humans{
		Id:       m.ID.Hex(),
		Name:     m.Name,
		LastName: m.LastName,
		Age:      m.Age,
		Email:    m.Email,
		Tel:      m.Tel,
	}
}

func NewHumanRepositoryMongo(client *mongo.Client, dbName string) *humanRepositoryMongo {
	return &humanRepositoryMongo{col: client.Database(dbName).Collection("humans")}
}

func (r *humanRepositoryMongo) GetPeoples() ([]entity.Humans, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []entity.Humans
	for cur.Next(ctx) {
		var m HumanMongo
		if err := cur.Decode(&m); err != nil {
			return nil, err
		}
		out = append(out, humanMongoToEn(m))
	}
	return out, nil
}

func (r *humanRepositoryMongo) GetPerson(id string) (*entity.Humans, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var m HumanMongo
	if err := r.col.FindOne(ctx, bson.M{"id": id}).Decode(&m); err != nil {
		return nil, err
	}
	en := humanMongoToEn(m)
	return &en, nil
}

func (r *humanRepositoryMongo) AddPerson(p entity.Humans) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if p.Id == "" {
		p.Id = primitive.NewObjectID().Hex()
	}

	m := humanEnToMongo(p)
	filter := bson.M{"id": m.ID}
	update := bson.M{"$setOnInsert": m}
	opts := options.Update().SetUpsert(true)

	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	return err
}
