package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hexagonal2/core/middleware"
	"hexagonal2/core/entity"
	e "hexagonal2/pkg/errors"
	"strconv"
)

type userRepositoryMongo struct {
	col *mongo.Collection
}

type UserMongo struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID   string `bson:"user_id" json:"user_id"`
	Name     string `bson:"name" json:"name"`
	LastName string `bson:"last_name" json:"last_name"`
	Age      int    `bson:"age" json:"age"`
	Email    string `bson:"email" json:"email"`
	Tel      string `bson:"tel" json:"tel"`
	Password string `bson:"password" json:"password"`
	Role     string `bson:"role" json:"role"` // Admin, User
}

func userEnToMongo(u entity.User) UserMongo {
	return UserMongo{
		ID:       primitive.NewObjectID(),
		UserID:   u.UserID,
		Name:     u.Name,
		LastName: u.LastName,
		Age:      u.Age,
		Email:    u.Email,
		Tel:      u.Tel,
		Password: u.Password,
		Role:     u.Role,
	}
}

func userMongoToEn(m UserMongo) entity.User {
	return entity.User{
		Id:       m.ID.Hex(),
		UserID:   m.UserID,
		Name:     m.Name,
		LastName: m.LastName,
		Age:      m.Age,
		Email:    m.Email,
		Tel:      m.Tel,
		Password: m.Password,
		Role:     m.Role,
	}
}

func NewUserRepositoryMongo(client *mongo.Client, dbName string) *userRepositoryMongo {
	return &userRepositoryMongo{col: client.Database(dbName).Collection("users")}
}

func (r *userRepositoryMongo) GetUsers() ([]entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil,e.ErrInternalServer 
	}
	defer cur.Close(ctx)

	var out []entity.User
	for cur.Next(ctx) {
		var m UserMongo
		if err := cur.Decode(&m); err != nil {
			return nil, e.ErrInternalServer
		}
		out = append(out, userMongoToEn(m))
	}
	return out, nil
}

func (r *userRepositoryMongo) GetUser(id string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var m UserMongo
	if err := r.col.FindOne(ctx, bson.M{"id": id}).Decode(&m); err != nil {
		return nil, e.ErrUserNotFound
	}
	en := userMongoToEn(m)
	return &en, nil
}

func (r *userRepositoryMongo) AddUser(p entity.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if p.Id == "" {
		p.Id = primitive.NewObjectID().Hex()
	}
	if p.UserID == "" {
		p.UserID = r.GenerateAccountID()
	}

	Password, err := middleware.HashPassword(p.Password)
	if err != nil {
		return e.ErrInternalServer
	}
	p.Password = Password


	_, err = r.col.InsertOne(ctx, userEnToMongo(p))
	return err
}

func (r *userRepositoryMongo) GetUserByEmail(email string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var m UserMongo
	if err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&m); err != nil {
		return nil, e.ErrUserNotFound
	}
	en := userMongoToEn(m)
	return &en, nil
}

func (r *userRepositoryMongo) GenerateAccountID() string {	
	var user entity.User
	seq := 477631 // default start
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()	
	err := r.col.FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.M{"user_id": -1})).Decode(&user)
	if err == nil && user.UserID != "" {
		if n, parseErr := strconv.Atoi(user.UserID); parseErr == nil {
			seq = n + 1
		}
	}
	return strconv.Itoa(seq)
}