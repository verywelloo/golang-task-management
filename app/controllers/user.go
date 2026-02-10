package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	res "github.com/verywelloo/3-go-echo-task-management/app/dto/response"
	m "github.com/verywelloo/3-go-echo-task-management/app/models"
	s "github.com/verywelloo/3-go-echo-task-management/app/services"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := s.AppInstance.Collections.Users

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "failed to get users",
		})
	}

	var users m.User
	if err := cursor.All(ctx, &users); err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "failed to decode users",
			Details: err.Error(),
		})
	}

	defer cursor.Close(ctx)

	return c.JSON(http.StatusOK, res.Result{
		Status:  http.StatusOK,
		Message: "successfully to get users",
	})
}
