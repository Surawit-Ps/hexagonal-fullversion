package handler

import (
	"hgo/core/entity"
	"hgo/core/ports"

	"github.com/gofiber/fiber/v2"
)

type humanHandler struct {
	humanSer ports.HumanServices
}

func NewHumanHandler(humanSer ports.HumanServices) humanHandler {
	return humanHandler{humanSer: humanSer}
}

func (h humanHandler) GetAllUsers(c *fiber.Ctx) error {
	humans, err := h.humanSer.GetAllUser()
	if err != nil {
		return c.JSON(fiber.ErrInternalServerError)
	}
	return c.JSON(humans)
}

func (h humanHandler) GetAUser(c *fiber.Ctx) error {
	human, err := h.humanSer.GetUser(c.Params("id"))
	if err != nil {
		return c.JSON(fiber.ErrNotFound)
	}
	return c.JSON(human)
}

func (h humanHandler) AddUser(c *fiber.Ctx) error {
	var person entity.Humans
	if err := c.BodyParser(&person); err != nil {
		return c.JSON(fiber.ErrInternalServerError)
	}
	if err := h.humanSer.AddUser(person); err != nil {
		return c.JSON(fiber.ErrInternalServerError)
	}
	return c.JSON(c.Status(200))
}
