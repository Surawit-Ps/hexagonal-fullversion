package handler

import (
	"hexagonal2/core/entity"
	"hexagonal2/core/ports"
	"hexagonal2/pkg/logs"
	"github.com/gofiber/fiber/v2"
)

type bookingHandler struct{
	bookingSer ports.BookingServices
	logger *logs.ZapLogger
}

func NewBookingHandler(bookingSer ports.BookingServices, logger *logs.ZapLogger) bookingHandler{
	return bookingHandler{bookingSer: bookingSer, logger: logger}
}

func(h bookingHandler)GetAllBookings(c *fiber.Ctx)error{
	booking ,err:= h.bookingSer.GetAllBookings()
	if  err != nil{
		return handleError(c,err)
	}
	h.logger.Info("Fetch all bookings successfully")
	return c.JSON(newResponseSuccess("Fetch all bookings successfully", booking))
}

func(h bookingHandler)GetABooking(c *fiber.Ctx)error{
	booking,err:=h.bookingSer.GetBooking(c.Params("id"))
	if err!=nil{
		return handleError(c,err)
	}
	h.logger.Info("Fetch booking with ID: " + c.Params("id"))
	return c.Status(fiber.StatusOK).JSON(newResponse(true, "Success", booking))
}

func(h bookingHandler)AddBooking(c *fiber.Ctx)error{
	var booking entity.Booking
	err := c.BodyParser(&booking)

	if err != nil{
		return handleError(c,err)
	}
	if err := h.bookingSer.AddBooking(booking); err != nil{
		return handleError(c, err)
	}

	h.logger.Info("Booking added with ID: " + booking.BookingID)
	return c.Status(fiber.StatusCreated).JSON(newResponse(true, "Success", nil))
}

func(h bookingHandler)UpdateBooking(c *fiber.Ctx)error{
	var booking entity.Booking
	err := c.BodyParser(&booking)

	if err != nil{
		return handleError(c,err)
	}	
	if err := h.bookingSer.UpdateBooking(c.Params("id"), booking); err != nil{
		return handleError(c, err)
	}
	h.logger.Info("Booking updated with ID: " + c.Params("id"))
	return c.Status(fiber.StatusOK).JSON(newResponse(true, "Success", nil))
}

func(h bookingHandler)DeleteBooking(c *fiber.Ctx)error{
	if err := h.bookingSer.DeleteBooking(c.Params("id")); err != nil{
		return handleError(c, err)
	}
	h.logger.Info("Booking deleted with ID: " + c.Params("id"))
	return c.Status(fiber.StatusOK).JSON(newResponse(true, "Success", nil))
}

