package main

import (
	"context"
	"log"
	"time"

	"hgo/adapter/handler"
	"hgo/adapter/repository"
	"hgo/core/entity"
	"hgo/core/ports"
	"hgo/core/service"

	"github.com/gofiber/fiber/v2"
	// "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("hgo.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&repository.HumanDB{}, &repository.DogsModel{}); err != nil {
		log.Fatal(err)
	}

	// repositories (default to GORM sqlite implementations)
	var humanRepo ports.HumanRepository = repository.NewHumanReposityDB(db)
	var dogRepo ports.DogsRepository = repository.NewDogsRepositoryDB(db)

	// try connecting to MongoDB and swap repositories if available
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err == nil {
		if err = mongoClient.Ping(ctx, nil); err == nil {
			log.Println("connected to mongodb, switching repositories")
			humanRepo = repository.NewHumanRepositoryMongo(mongoClient, "hgo")
			dogRepo = repository.NewDogsRepositoryMongo(mongoClient, "hgo")
		} else {
			log.Println("mongo ping failed:", err)
		}
	} else {
		log.Println("mongo connect failed:", err)
	}

	// services
	humanSrv := service.NewHumanServive(humanRepo)
	dogSrv := service.NewDogService(dogRepo)

	// handlers
	dh := handler.NewDogHandler(dogSrv)
	hh := handler.NewHumanHandler(humanSrv)

	// Seed mock data via repository interfaces (works for both GORM and Mongo)
	people := []entity.Humans{
		{Name: "John", LastName: "Doe", Age: 30, Email: "john@example.com", Tel: "123-456"},
		{Name: "Jane", LastName: "Smith", Age: 28, Email: "jane@example.com", Tel: "987-654"},
	}
	for _, p := range people {
		if err := humanRepo.AddPerson(p); err != nil {
			log.Println("seed human failed:", err)
		}
	}

	humans, err := humanRepo.GetPeoples()
	if err != nil {
		log.Println("failed to list humans after seed:", err)
	}
	emap := map[string]string{}
	for _, h := range humans {
		emap[h.Email] = h.Id
	}

	dogs := []entity.Dogs{
		{Name: "Rex", Age: 3, Colour: "Brown", HumanID: emap["john@example.com"]},
		{Name: "Milo", Age: 2, Colour: "Black", HumanID: emap["jane@example.com"]},
	}
	for _, d := range dogs {
		if err := dogRepo.AddDog(d, d.HumanID); err != nil {
			log.Println("seed dog failed:", err)
		}
	}

	// Seed idempotent mock data for testing using FirstOrCreate
	// john := repository.HumanDB{Id: uuid.New().String(), Name: "John", LastName: "Doe", Age: 30, Email: "john@example.com", Tel: "123-456"}
	// jane := repository.HumanDB{Id: uuid.New().String(), Name: "Jane", LastName: "Smith", Age: 28, Email: "jane@example.com", Tel: "987-654"}

	// var johnTmp, janeTmp repository.HumanDB
	// if err := db.FirstOrCreate(&johnTmp, repository.HumanDB{Id: john.Id, Email: john.Email, Name: john.Name, LastName: john.LastName, Age: john.Age, Tel: john.Tel}).Error; err != nil {
	// 	log.Println("seed human john:", err)
	// }
	// if err := db.FirstOrCreate(&janeTmp, repository.HumanDB{Id: jane.Id, Email: jane.Email, Name: jane.Name, LastName: jane.LastName, Age: jane.Age, Tel: jane.Tel}).Error; err != nil {
	// 	log.Println("seed human jane:", err)
	// }

	// var johnDB, janeDB repository.HumanDB
	// db.First(&johnDB, "id = ?", john.Id)
	// db.First(&janeDB, "id = ?", jane.Id)

	// dog1 := repository.DogsModel{Id: uuid.New().String(), Name: "Rex", Age: 3, Colour: "Brown", HumanID: johnDB.Id}
	// if err := db.FirstOrCreate(&dog1, repository.DogsModel{Name: dog1.Name, HumanID: dog1.HumanID}).Error; err != nil {
	// 	log.Println("seed dog rex:", err)
	// }

	// dog2 := repository.DogsModel{Id: uuid.New().String(), Name: "Milo", Age: 2, Colour: "Black", HumanID: janeDB.Id}
	// if err := db.FirstOrCreate(&dog2, repository.DogsModel{Name: dog2.Name, HumanID: dog2.HumanID}).Error; err != nil {
	// 	log.Println("seed dog milo:", err)
	// }

	app := fiber.New()

	app.Get("/dogs", dh.GetAllDogs)
	app.Get("/dogs/:id", dh.GetADogs)
	app.Post("/dogs", dh.AddDog)

	app.Get("/humans", hh.GetAllUsers)
	app.Get("/humans/:id", hh.GetAUser)
	app.Post("/humans", hh.AddUser)

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
