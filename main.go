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


	userRepo,petRepo,bookingRepo, err := connectDatabase(true) // pass false to attempt MongoDB connection
	if err != nil {
		log.Fatal("failed to connect to any database:", err)
	}
	logger, err := logs.NewZapLogger()
	if err != nil {
		log.Fatal("failed to initialize logger:", err)
	}

	// services
	userSrv := service.NewUserService(userRepo, *logger)
	petSrv := service.NewPetService(petRepo)
	bookingSrv := service.NewBookingService(bookingRepo)
	// dogSrv := service.NewDogService(dogRepo)

	// handlers
	// dh := handler.NewDogHandler(dogSrv, logger)
	hh := handler.NewUserHandler(userSrv, logger)
	ph := handler.NewPetHandler(petSrv, logger)
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
	app.Get("/pets", ph.GetAllPets)
	app.Get("/pets/:id", ph.GetAPet)
	app.Post("/pets", ph.AddPet)
	app.Get("/bookings", bh.GetAllBookings)
	app.Get("/bookings/:id", bh.GetABooking)
	app.Post("/bookings", bh.AddBooking)
	app.Put("/bookings/:id", bh.UpdateBooking)
	app.Delete("/bookings/:id", bh.DeleteBooking)

	// pet := []entity.Pet{
	// 	{Name: "Buddy", Species: "Dog", Breed: "Labrador", Age: 3, Weight: 20.5},
	// 	{Name: "Mittens", Species: "Cat", Breed: "Siamese", Age: 2, Weight: 5.0},
	// }
	// for _, p := range pet {
	// 	if err := petRepo.AddPet(p, "1"); err != nil {
	// 		log.Println("failed to add pet:", err)
	// 	}
	// }



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


func connectDatabase(flag bool) (userRepo ports.UserRepository,petRepo ports.PetRepository,bookingRepo ports.BookingRepository, err error) {
	if flag {
		db, err := gorm.Open(sqlite.Open("hgo.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&sql.User{}, &sql.Booking{}, &sql.Pet{}); err != nil {
		log.Fatal(err)
	}

	userRepo = sql.NewUserRepositoryDB(db)
	// dogRepo = sql.NewDogsRepositoryDB(db)
	petRepo = sql.NewPetRepositoryDB(db)
	bookingRepo = sql.NewBookingRepositoryDB(db)

	return userRepo, petRepo, bookingRepo, nil


	}else{
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		mongoClient, err := mg.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
		if err == nil {
		if err = mongoClient.Ping(ctx, nil); err == nil {
			log.Println("connected to mongodb, switching repositories")
			userRepo = mongo.NewUserRepositoryMongo(mongoClient, "hgo")
			// dogRepo = mongo.NewDogsRepositoryMongo(mongoClient, "hgo")
			petRepo = mongo.NewPetRepositoryMongo(mongoClient, "hgo")
			bookingRepo = mongo.NewBookingRepositoryMongo(mongoClient, "hgo")
			return userRepo, petRepo, bookingRepo, nil
		} else {
			log.Println("mongo ping failed:", err)
		}
	} else {
		log.Println("mongo connect failed:", err)
	}

	}
	
	// try connecting to MongoDB and swap repositories if available

	return nil, nil, nil, err
}

