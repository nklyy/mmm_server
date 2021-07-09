package routes

import (
	"mmm_server/controllers"

	"github.com/gofiber/fiber/v2"
)

type Route struct {
	controll *controllers.Controller
}

func NewRoute(controll *controllers.Controller) *Route {
	return &Route{controll: controll}
}

func (r *Route) Initialroute(route fiber.Router) {
	user := route.Group("/user")

	{
		user.Get("/", r.controll.GetAllUsers)
	}
}
