package unused 


// import (
// 	"hexagonal2/core/entity"
// 	"hexagonal2/core/ports"
// 	"hexagonal2/pkg/logs"
// 	"github.com/gofiber/fiber/v2"
// )

// type dogHandler struct{
// 	dogSer ports.DogServices
// 	logger *logs.ZapLogger
// }

// func NewDogHandler(dogSer ports.DogServices, logger *logs.ZapLogger) dogHandler{
// 	return dogHandler{dogSer: dogSer, logger: logger}
// }

// func(h dogHandler)GetAllDogs(c *fiber.Ctx)error{
// 	dog ,err:= h.dogSer.GetAllDogs()
// 	if  err != nil{
// 		return handleError(c,err)
// 	}
// 	h.logger.Info("Fetch all dogs successfully")
// 	return c.JSON(newResponseSuccess("Fetch all dogs successfully", dog))
// }

// func(h dogHandler)GetADogs(c *fiber.Ctx)error{
// 	dog,err:=h.dogSer.GetDog(c.Params("id"))
// 	if err!=nil{
// 		return handleError(c,err)
// 	}
// 	h.logger.Info("Fetch dog with ID: " + c.Params("id"))
// 	return c.Status(fiber.StatusOK).JSON(newResponse(true, "Success", dog))
// }

// func(h dogHandler)AddDog(c *fiber.Ctx)error{
// 	var dog entity.Dogs
// 	err := c.BodyParser(&dog)
// 	if err != nil{
// 		return handleError(c,err)
// 	}
// 	id := dog.UserID
// 	if err := h.dogSer.AddDog(dog, id); err != nil{
// 		return handleError(c, err)
// 	}
// 	h.logger.Info("Dog added with ID: " + dog.Id)
// 	return c.Status(fiber.StatusCreated).JSON(newResponse(true, "Success", nil))
// }