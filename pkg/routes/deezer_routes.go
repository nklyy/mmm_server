package routes

import "github.com/gofiber/fiber/v2"

func DeezerRoutes(app *fiber.App) {
	api := app.Group("/v1")

	api.Get("/deezer", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World 👋!")
	})

	//app.Post("/", func(ctx *fiber.Ctx) error {
	//	var d Data
	//
	//	err := ctx.BodyParser(&d)
	//
	//	if err != nil {
	//		return fiber.NewError(fiber.StatusBadRequest, "Can't parse body!")
	//	}
	//
	//	//result, err := json.Marshal(`{"id": 42, "username": "rvasily", "phone": "123"}`)
	//	//if err != nil {
	//	//	panic(err)
	//	//}
	//	//
	//	//return ctx.Send(result)
	//
	//	return ctx.SendString("OK")
	//})
}
