package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	res "github.com/verywelloo/3-go-echo-task-management/app/dto/response"
	m "github.com/verywelloo/3-go-echo-task-management/app/models"
	s "github.com/verywelloo/3-go-echo-task-management/app/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTask(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	taskCollection := s.AppInstance.Collections.Tasks

	projectIDStr := c.Param("project_id")

	if err := c.Bind(&)

	projectID, err := primitive.ObjectIDFromHex(projectIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, res.Result{
			Status:  http.StatusBadRequest,
			Message: "invalid project id in param",
			Details: err.Error(),
		})
	}

	insert := m.Task{
		ID:        primitive.NewObjectID(),
		ProjectID: projectID,
		//StartDate:
		//EndDateDate:
	}

	if _, err := taskCollection.InsertOne(ctx, insert); err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "failed to create a task",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res.Result{
		Status:  http.StatusOK,
		Message: "successfully create a task",
	})
}
