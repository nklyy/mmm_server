package controllers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"mmm_server/platform/database"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Text      string             `bson:"text"`
	Age       int                `bson:"age"`
}

func GetAllUsers(ctx *fiber.Ctx) error {
	db, err := database.MongoDbConnection()

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "books were not found",
			"count": 0,
			"books": nil,
		})
	}

	cursor, err := db.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var users []bson.M
	err = cursor.All(context.TODO(), &users)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(users)

	return nil
}
