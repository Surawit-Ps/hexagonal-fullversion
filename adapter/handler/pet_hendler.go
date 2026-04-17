package handler

import (
	"hexagonal2/core/entity"
	"hexagonal2/core/ports"
	"hexagonal2/pkg/logs"

	"github.com/gofiber/fiber/v2"
)

type petHandler struct{
	petSer ports.PetServices
	logger *logs.ZapLogger
}

func NewPetHandler(petSer ports.PetServices, logger *logs.ZapLogger) petHandler{
	return petHandler{petSer: petSer, logger: logger}
}

func(h petHandler)GetAllPets(c *fiber.Ctx)error{
	pet ,err:= h.petSer.GetAllPets()
	if  err != nil{
		return handleError(c,err)
	}
	h.logger.Info("Fetch all pets successfully")
	return c.JSON(newResponseSuccess("Fetch all pets successfully", pet))
}

func(h petHandler)GetAPet(c *fiber.Ctx)error{
	pet,err:=h.petSer.GetPet(c.Params("id"))
	if err!=nil{
		return handleError(c,err)
	}
	h.logger.Info("Fetch pet with ID: " + c.Params("id"))
	return c.Status(fiber.StatusOK).JSON(newResponse(true, "Success", pet))
}

func(h petHandler)AddPet(c *fiber.Ctx)error{
	var pet entity.Pet
	err := c.BodyParser(&pet)
	if err != nil{
		return handleError(c,err)
	}
	id := pet.OwnerID
	if err := h.petSer.AddPet(pet, id); err != nil{
		return handleError(c, err)
	}
	h.logger.Info("Pet added with ID: " + pet.ID)
	return c.Status(fiber.StatusCreated).JSON(newResponse(true, "Success", nil))
}

func(h petHandler)UpdatePet(c *fiber.Ctx)error{
	var pet entity.Pet
	err := c.BodyParser(&pet)
	if err != nil{
		return handleError(c,err)
	}
	if err := h.petSer.UpdatePet(c.Params("id"), pet); err != nil{
		return handleError(c, err)
	}
	h.logger.Info("Pet updated with ID: " + c.Params("id"))
	return c.Status(fiber.StatusOK).JSON(newResponse(true, "Success", nil))
}

func(h petHandler)DeletePet(c *fiber.Ctx)error{
	if err := h.petSer.DeletePet(c.Params("id")); err != nil{
		return handleError(c, err)
	}
	h.logger.Info("Pet deleted with ID: " + c.Params("id"))
	return c.Status(fiber.StatusOK).JSON(newResponse(true, "Success", nil))
}

