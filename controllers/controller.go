package controllers

import (
	"github.com/gofiber/fiber/v2"
	"mmm_server/repositories"
)

type User interface {
	GetAllUsers(*fiber.Ctx) error
}

type Controller struct {
	User
}

func NewController(repos *repositories.Repository) *Controller {
	return &Controller{
		User: NewUserController(repos.User),
	}
}
