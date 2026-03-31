package handler

import (
	"hexagonal2/core/entity"
	"hexagonal2/core/ports"
	"hexagonal2/pkg/logs"
	"github.com/gofiber/fiber/v2"
	"hexagonal2/core/middleware"
)

type userHandler struct {
	userSer ports.UserServices
	logger *logs.ZapLogger
}

func NewUserHandler(userSer ports.UserServices, logger *logs.ZapLogger) userHandler {
	return userHandler{userSer: userSer, logger: logger}
}

func (h userHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userSer.GetAllUser()
	if err != nil {
		return handleError(c,err)
	}
	return c.JSON(newResponse(true, "Success", users))
}

func (h userHandler) GetAUser(c *fiber.Ctx) error {
	user, err := h.userSer.GetUser(c.Params("id"))
	if err != nil {
		return handleError(c, err)
	}
	h.logger.Info("Get user with ID: " + c.Params("id"))
	return c.JSON(newResponse(true, "Success", user))
}

func (h userHandler) AddUser(c *fiber.Ctx) error {
	var person entity.User
	if err := c.BodyParser(&person); err != nil {
		return handleError(c, err)
	}
	if err := h.userSer.AddUser(person); err != nil {
		return handleError(c, err)
	}
	h.logger.Info("User added with ID: " + person.Id)
	return c.JSON(newResponse(true, "Success", nil))
}

func (h userHandler) Login(c *fiber.Ctx) error {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return handleError(c, err)
	}
	token, err := h.userSer.Login(req.Email, req.Password)
	if err != nil {
		return handleError(c, err)
	}
	middleware.SetCookies(c, token)
	h.logger.Info("User logged in with ID: " + c.Params("id"))
	return c.JSON(newResponse(true, "Login successful", newAuthResponse(token)))
}