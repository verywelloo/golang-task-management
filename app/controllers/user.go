package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	s "github.com/verywelloo/3-go-echo-task-management/app/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := s.AppInstance.Collections.Users

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, "user not found")
		} else {
			return c.JSON(http.StatusInternalServerError, "error")
		}
	}

	defer cursor.Close(ctx)

	return c.JSON(http.StatusOK, cursor)
}
