package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	req "github.com/verywelloo/3-go-echo-task-management/app/dto/request"
	m "github.com/verywelloo/3-go-echo-task-management/app/models"

	"go.mongodb.org/mongo-driver/bson"
)

func Register(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var payload req.RegisterPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	var user m.User
	if err := UserCollection.FindOne(ctx, bson.M{"email": payload.Email}).Decode(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, "error in finding user with email")
	}

	if user.Name == "" {
		insert := m.User{
			Email:     payload.Email,
			Name:      payload.Name,
			Password:  payload.Password,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := UserCollection.InsertOne(ctx, insert)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "error in inserting user")
		}
	}

	return c.JSON(http.StatusOK, "OK")
}
