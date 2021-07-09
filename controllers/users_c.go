package controllers

import (
	"mmm_server/repositories"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	repo repositories.User
}

func NewUserController(repo repositories.User) *UserController {
	return &UserController{
		repo: repo,
	}
}

func (uc *UserController) GetAllUsers(ctx *fiber.Ctx) error {

	users := uc.repo.GetAllUsers()

	return ctx.JSON(users)
}
