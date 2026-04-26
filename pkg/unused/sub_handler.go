package unused

// import (
// 	"hexagonal2/core/entity"
// 	"hexagonal2/core/ports"
// 	"hexagonal2/pkg/errors"
// 	"github.com/gofiber/fiber/v2"
// 	e "errors"
// )

// type subHandler struct{
// 	subSer ports.SubscriptionService
// }

// func NewSubHandler(subSer ports.SubscriptionService) subHandler{
// 	return subHandler{subSer: subSer}
// }

// func(h subHandler)CreateSubscription(c *fiber.Ctx)error{
// 	var sub entity.Subscription
// 	err := c.BodyParser(&sub)
// 	if err != nil{
// 		return handleError(c, err)
// 	}
// 	if err := h.subSer.CreateSubscription(sub); err != nil{
// 		return handleError(c, err)
// 	}
// 	return c.Status(fiber.StatusCreated).JSON(newResponse(true, "Subscription created successfully", nil))
// }

// func(h subHandler)GetSubscriptionByUserID(c *fiber.Ctx)error{
// 	userID := c.Params("userID")
// 	sub, err := h.subSer.GetSubscriptionByUserID(userID)	
// 	if err != nil{
// 		if e.Is(err, errors.ErrNotFound) {
// 			return handleError(c, errors.ErrNotFound)
// 		}
// 		return handleError(c, err)
// 	}
// 	return c.Status(fiber.StatusOK).JSON(newResponse(true, "Subscription fetched successfully", sub))
// }

