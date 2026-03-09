package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	req "github.com/verywelloo/3-go-echo-task-management/app/dto/request"
	res "github.com/verywelloo/3-go-echo-task-management/app/dto/response"
	m "github.com/verywelloo/3-go-echo-task-management/app/models"
	s "github.com/verywelloo/3-go-echo-task-management/app/services"
)

func CreateProject(c echo.Context) error {
	userCollection := s.AppInstance.Collections.Users
	projectCollection := s.AppInstance.Collections.Projects
	ctx := c.Request().Context()

	// session, err := s.GetSessionCache(c)
	// if err != nil {
	// 	return c.JSON(http.StatusUnauthorized, res.Result{
	// 		Status: http.StatusUnauthorized,
	// 		Message: "unauthorized",
	// 		Details: err.Error(),
	// 	})
	// }

	var payload req.CreateProjectPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "failed to parse payload",
			Details: err.Error(),
		})
	}

	if payload.StartDate != "" {

	}

	if payload.EndDate != "" {

	}

	newProject := m.Project{
		ID:   primitive.NewObjectID(),
		Name: payload.Name,
		//StartDate
		//EndDate
	}

	_, err := projectCollection.InsertOne(ctx, newProject)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "failed to create a project",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res.Result{
		Status:  http.StatusOK,
		Message: "successfully create a project",
	})
}

func GetProject(c echo.Context) error {
	userCollection := s.AppInstance.Collections.Users
	ctx := c.Request().Context()

	session, err := s.GetSessionCache(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, res.Result{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
			Details: err.Error(),
		})
	}

	var projectPermission m.ProjectPermission
	if err := userCollection.FindOne(ctx, bson.M{
		"user_id":    session.ID,
		"deleted_at": nil,
	}).Decode(&projectPermission); err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "failed to retrieve project-permission",
			Details: err.Error(),
		})
	}

	return nil
}
