package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	req "github.com/verywelloo/3-go-echo-task-management/app/dto/request"
	m "github.com/verywelloo/3-go-echo-task-management/app/models"
	s "github.com/verywelloo/3-go-echo-task-management/app/services"

	"go.mongodb.org/mongo-driver/bson"
)

func Register(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := s.AppInstance.Collections.Users

	var payload req.RegisterPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request payload")
	}

	// check exists email
	var user m.User
	if err := userCollection.FindOne(ctx, bson.M{"email": payload.Email}).Decode(&user); err == nil {
		return c.JSON(http.StatusBadRequest, "user already exists")
	}

	// new user, create one
	if user.Name == "" {
		insert := m.User{
			Email:     payload.Email,
			Name:      payload.Name,
			Password:  payload.Password,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := userCollection.InsertOne(ctx, insert)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "error in inserting user")
		}
	}

	return c.JSON(http.StatusOK, "successfully to create a user")
}
