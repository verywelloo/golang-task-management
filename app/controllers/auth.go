package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	req "github.com/verywelloo/3-go-echo-task-management/app/dto/request"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var payload req.RegisterPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	email, err := UserCollection.Find(ctx, bson.M{"email": payload.Email})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, "email not found")
		}
		return c.JSON(http.StatusInternalServerError, "error in finding user with email")
	}
	defer email.Close(ctx)

	// userData := m.User{
	// 	Email: payload.Email,
	// Name: payload.Name,
	// Password: ,//payload.Password,  hashing
	// CreatedAt: time.Now(),
	// UpdatedAt: time.Now(),
	// }

	// user, err := UserCollection.InsertOne(ctx, user)

	return c.JSON(http.StatusOK, "OK")
}
