package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	req "github.com/verywelloo/3-go-echo-task-management/app/dto/request"
	res "github.com/verywelloo/3-go-echo-task-management/app/dto/response"
	m "github.com/verywelloo/3-go-echo-task-management/app/models"
	s "github.com/verywelloo/3-go-echo-task-management/app/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTask(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	taskCollection := s.AppInstance.Collections.Tasks

	var payload req.CreateTask
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, res.Result{
			Status:  http.StatusBadRequest,
			Message: "invalid payload in creating a task",
			Details: err.Error(),
		})
	}

	projectID, err := primitive.ObjectIDFromHex(payload.ProjectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, res.Result{
			Status:  http.StatusBadRequest,
			Message: "invalid project id in param",
			Details: err.Error(),
		})
	}

	loc, err := time.LoadLocation("Asia/Bangkok")
	startDate, err := time.ParseInLocation(time.DateOnly, payload.StartDate, loc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "failed to parse start date",
			Details: err.Error(),
		})
	}

	endDate, err := time.ParseInLocation(time.DateOnly, payload.EndDate, loc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "failed to parse end date",
			Details: err.Error(),
		})
	}

	insert := m.Task{
		ID:        primitive.NewObjectID(),
		ProjectID: projectID,
		StartDate: startDate,
		EndDate:   endDate,
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

func GetTasks(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	taskCollection := s.AppInstance.Collections.Tasks

	cur, err := taskCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "failed to retrieve tasks",
			Details: err.Error(),
		})
	}

	var response res.GetTask

	return nil
}
