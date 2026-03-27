package handler

import (
	"hgo/core/entity"
	"hgo/core/ports"

	"github.com/gofiber/fiber/v2"
)

type dogHandler struct{
	dogSer ports.DogServices
}

func NewDogHandler(dogSer ports.DogServices)dogHandler{
	return dogHandler{dogSer: dogSer}
}

func(h dogHandler)GetAllDogs(c *fiber.Ctx)error{
	dog ,err:= h.dogSer.GetAllDogs()
	if  err != nil{
		return c.JSON(fiber.ErrInternalServerError)
	}
	return c.JSON(dog)
}

func(h dogHandler)GetADogs(c *fiber.Ctx)error{
	dog,err:=h.dogSer.GetDog(c.Params("id"))
	if err!=nil{
		return c.JSON(fiber.ErrNotFound)
	}
	return c.JSON(dog)
}

func(h dogHandler)AddDog(c *fiber.Ctx)error{
	var dog []entity.Dogs
	err := c.BodyParser(&dog)
	if err != nil{
		return c.JSON(fiber.ErrInternalServerError)
	}
	return c.JSON(c.Status(200))
}