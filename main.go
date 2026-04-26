package main

import (
	"context"
	"log"
	"time"

	"hexagonal2/adapter/handler"
	"hexagonal2/adapter/repository/sql"
	"hexagonal2/adapter/repository/mongo"

	// "hexagonal2/core/entity"
	"hexagonal2/core/ports"
	"hexagonal2/core/service"
	"hexagonal2/core/middleware"
	"github.com/gofiber/fiber/v2"
	// "github.com/google/uuid"
	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"hexagonal2/pkg/logs"
)

func main() {


	userRepo,bookingRepo, err := connectDatabase(true) // pass false to attempt MongoDB connection
	if err != nil {
		log.Fatal("failed to connect to any database:", err)
	}
	logger, err := logs.NewZapLogger()
	if err != nil {
		log.Fatal("failed to initialize logger:", err)
	}

	// services
	userSrv := service.NewUserService(userRepo, *logger)
	// petSrv := service.NewPetService(petRepo)
	bookingSrv := service.NewBookingService(bookingRepo)
	// dogSrv := service.NewDogService(dogRepo)

	// handlers
	// dh := handler.NewDogHandler(dogSrv, logger)
	hh := handler.NewUserHandler(userSrv, logger)
	// ph := handler.NewPetHandler(petSrv, logger)
	bh := handler.NewBookingHandler(bookingSrv, logger)


	app := fiber.New()
	middleware.CORS(app) // apply CORS middleware globally

	app.Post("/login",hh.Login)
	app.Use(middleware.Authorizes()) // apply authorization middleware to all routes below
	// app.Get("/dogs", dh.GetAllDogs)
	// app.Get("/dogs/:id", dh.GetADogs)
	// app.Post("/dogs", dh.AddDog)

	app.Get("/users", hh.GetAllUsers)
	app.Get("/users/:id", hh.GetAUser)
	app.Post("/users", hh.AddUser)
	// app.Get("/pets", ph.GetAllPets)
	// app.Get("/pets/:id", ph.GetAPet)
	// app.Post("/pets", ph.AddPet)
	app.Get("/bookings", bh.GetAllBookings)
	app.Get("/bookings/:id", bh.GetABooking)
	app.Post("/bookings", bh.AddBooking)
	app.Put("/bookings/:id", bh.UpdateBooking)
	app.Delete("/bookings/:id", bh.DeleteBooking)


	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}


func connectDatabase(flag bool) (userRepo ports.UserRepository,bookingRepo ports.BookingRepository, err error) {
	if flag {
		db, err := gorm.Open(sqlite.Open("hgo.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&sql.User{}, &sql.Booking{}); err != nil {
		log.Fatal(err)
	}

	userRepo = sql.NewUserRepositoryDB(db)
	// dogRepo = sql.NewDogsRepositoryDB(db)
	// petRepo = sql.NewPetRepositoryDB(db)
	bookingRepo = sql.NewBookingRepositoryDB(db)

	return userRepo, bookingRepo, nil


	}else{
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		mongoClient, err := mg.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
		if err == nil {
		if err = mongoClient.Ping(ctx, nil); err == nil {
			log.Println("connected to mongodb, switching repositories")
			userRepo = mongo.NewUserRepositoryMongo(mongoClient, "hgo")
			// dogRepo = mongo.NewDogsRepositoryMongo(mongoClient, "hgo")
			// petRepo = mongo.NewPetRepositoryMongo(mongoClient, "hgo")
			bookingRepo = mongo.NewBookingRepositoryMongo(mongoClient, "hgo")
			return userRepo, bookingRepo, nil
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

