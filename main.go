package main

import (
	"context"
	"log"
	"time"

	"hexagonal2/adapter/handler"
	"hexagonal2/adapter/repository"
	// "hexagonal2/core/entity"
	"hexagonal2/core/ports"
	"hexagonal2/core/service"
	"hexagonal2/core/middleware"
	"github.com/gofiber/fiber/v2"
	// "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"hexagonal2/pkg/logs"
)

func main() {


	userRepo, dogRepo, err := connectDatabase(true) // pass false to attempt MongoDB connection
	if err != nil {
		log.Fatal("failed to connect to any database:", err)
	}
	logger, err := logs.NewZapLogger()
	if err != nil {
		log.Fatal("failed to initialize logger:", err)
	}

	// services
	userSrv := service.NewUserService(userRepo)
	dogSrv := service.NewDogService(dogRepo)

	// handlers
	dh := handler.NewDogHandler(dogSrv, logger)
	hh := handler.NewUserHandler(userSrv, logger)


	app := fiber.New()
	middleware.CORS(app) // apply CORS middleware globally

	app.Post("/login",hh.Login)
	app.Use(middleware.Authorizes()) // apply authorization middleware to all routes below
	app.Get("/dogs", dh.GetAllDogs)
	app.Get("/dogs/:id", dh.GetADogs)
	app.Post("/dogs", dh.AddDog)

	app.Get("/users", hh.GetAllUsers)
	app.Get("/users/:id", hh.GetAUser)
	app.Post("/users", hh.AddUser)
	
	//use this block to seed initial data if needed
	// user := []entity.User{
	// 	{Name: "John", LastName: "Doe", Age: 30, Email: "john.doe@example.com", Tel: "1234567890", Password: "123", Role: "User"},
	// 	{Name: "Jane", LastName: "Smith", Age: 25, Email: "jane.smith@example.com", Tel: "0987654321", Password: "123", Role: "User"},
	// 	{Name: "Admin", LastName: "User", Age: 35, Email: "admin@example.com", Tel: "0000000000", Password: "123", Role: "Admin"},
	// }
	// for _, u := range user {
	// 	if err := userRepo.AddUser(u); err != nil {
	// 		log.Println("failed to add user:", err)
	// 	}
	// }

	// dog := []entity.Dogs{
	// 	{Name: "Buddy", Age: 3, Colour:"Black", UserID: "1"},
	// 	{Name: "Max", Age: 5, Colour:"White", UserID: "2"},
	// 	{Name: "Bella", Age: 2, Colour:"Brown", UserID: "1"},
	// }
	// for _, d := range dog {
	// 	if err := dogRepo.AddDog(d, d.UserID); err != nil {
	// 		log.Println("failed to add dog:", err)
	// 	}
	// }

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}


func connectDatabase(flag bool) (userRepo ports.UserRepository, dogRepo ports.DogsRepository,err error) {
	if flag {
		db, err := gorm.Open(sqlite.Open("hgo.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&repository.UserDB{}, &repository.DogsModel{}); err != nil {
		log.Fatal(err)
	}

	userRepo = repository.NewUserRepositoryDB(db)
	dogRepo = repository.NewDogsRepositoryDB(db)

	return userRepo, dogRepo, nil


	}else{
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
		if err == nil {
		if err = mongoClient.Ping(ctx, nil); err == nil {
			log.Println("connected to mongodb, switching repositories")
			userRepo = repository.NewUserRepositoryMongo(mongoClient, "hgo")
			dogRepo = repository.NewDogsRepositoryMongo(mongoClient, "hgo")
			return userRepo, dogRepo, nil
		} else {
			log.Println("mongo ping failed:", err)
		}
	} else {
		log.Println("mongo connect failed:", err)
	}

	}
	
	// try connecting to MongoDB and swap repositories if available

	return nil, nil, err
}

