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

type petRepositoryMongo struct {
	col *mongo.Collection
}

func NewPetRepositoryMongo(client *mongo.Client, dbName string) *petRepositoryMongo {
	return &petRepositoryMongo{col: client.Database(dbName).Collection("pets")}
}

type PetMongo struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    PetID     string `bson:"pet_id" json:"pet_id"`
    OwnerID   string `bson:"owner_id" json:"owner_id"`
    Name      string `bson:"name" json:"name"`
    Species   string `bson:"species" json:"species"` // dog, cat
    Breed     string `bson:"breed" json:"breed"`
    Age       int `bson:"age" json:"age"`
    Weight    float64 `bson:"weight" json:"weight"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

func petEnToMongo(p entity.Pet) PetMongo {
	return PetMongo{
		ID:        primitive.NewObjectID(),
		PetID:     p.PetID,
		OwnerID:   p.OwnerID,
		Name:      p.Name,
		Species:   p.Species,
		Breed:     p.Breed,
		Age:       p.Age,
		Weight:    p.Weight,
		CreatedAt: time.Now(),
	}
}

func petMongoToEn(m PetMongo) entity.Pet {
	return entity.Pet{
		ID:        m.ID.Hex(),
		PetID:     m.PetID,
		OwnerID:   m.OwnerID,
		Name:      m.Name,
		Species:   m.Species,
		Breed:     m.Breed,
		Age:       m.Age,
		Weight:    m.Weight,
		CreatedAt: m.CreatedAt,
	}
}

func (r *petRepositoryMongo) AddPet(p entity.Pet,h string) error {
	p.PetID = primitive.NewObjectID().Hex()
	p.OwnerID = h
	petMongo := petEnToMongo(p)
	_, err := r.col.InsertOne(context.Background(), petMongo)
	if err != nil {
		return e.ErrInternalServer
	}
	return nil
}

func (r *petRepositoryMongo) GetPets() ([]entity.Pet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, e.ErrInternalServer
	}
	defer cur.Close(ctx)

	var pets []entity.Pet
	for cur.Next(ctx) {
		var petMongo PetMongo
		if err := cur.Decode(&petMongo); err != nil {
			return nil, e.ErrInternalServer
		}
		pets = append(pets, petMongoToEn(petMongo))
	}
	return pets, nil
}

func (r *petRepositoryMongo) GetAPet(id string) (*entity.Pet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, e.ErrInvalidID
	}
	var petMongo PetMongo
	err = r.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&petMongo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, e.ErrNotFound
		}
		return nil, e.ErrInternalServer
	}
	pet := petMongoToEn(petMongo)
	return &pet, nil
}


func (r *petRepositoryMongo) UpdatePet(id string, p entity.Pet) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return e.ErrInvalidID
	}
	update := bson.M{
		"$set": bson.M{
			"name":    p.Name,
			"species": p.Species,
			"breed":   p.Breed,
			"age":     p.Age,
			"weight":  p.Weight,
		},
	}
	result, err := r.col.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return e.ErrInternalServer
	}
	if result.MatchedCount == 0 {
		return e.ErrNotFound
	}
	return nil
}

func (r *petRepositoryMongo) DeletePet(id string) error {
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