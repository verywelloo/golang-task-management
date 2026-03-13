package controllers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	req "github.com/verywelloo/3-go-echo-task-management/app/dto/request"
	res "github.com/verywelloo/3-go-echo-task-management/app/dto/response"
	m "github.com/verywelloo/3-go-echo-task-management/app/models"
	s "github.com/verywelloo/3-go-echo-task-management/app/services"
)

var (
	projectCollection           = s.AppInstance.Collections.Projects
	projectPermissionCollection = s.AppInstance.Collections.ProjectPermission
	userCollection              = s.AppInstance.Collections.Users
)

func CreateProject(c echo.Context) error {
	ctx := c.Request().Context()

	session, err := s.GetSessionCache(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, res.Result{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
			Details: err.Error(),
		})
	}

	userID, err := primitive.ObjectIDFromHex(session.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "invalid objectID",
			Details: err.Error(),
		})
	}

	var payload req.CreateProjectPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "failed to parse payload",
			Details: err.Error(),
		})
	}

	var startDate, endDate time.Time
	if payload.StartDate != "" {
		startDate, err = time.ParseInLocation("2006-01-02", payload.StartDate, time.Local)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, res.Result{
				Status:  http.StatusInternalServerError,
				Message: "failed to parse a start date",
				Details: err.Error(),
			})
		}
	}

	if payload.EndDate != "" {
		endDate, err = time.ParseInLocation("2006-01-02", payload.EndDate, time.Local)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, res.Result{
				Status:  http.StatusInternalServerError,
				Message: "failed to parse a end date",
				Details: err.Error(),
			})
		}
	}

	newProject := m.Project{
		ID:        primitive.NewObjectID(),
		Name:      payload.Name,
		StartDate: startDate,
		EndDate:   endDate,
	}

	result, err := projectCollection.InsertOne(ctx, newProject)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "failed to create a project",
			Details: err.Error(),
		})
	}

	// create permission
	projectPer := m.ProjectPermission{
		ID:        primitive.NewObjectID(),
		ProjectID: result.InsertedID.(primitive.ObjectID),
		UserID:    userID,
	}
	_, err = projectPermissionCollection.InsertOne(ctx, projectPer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "filed to create a project-permission",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res.Result{
		Status:  http.StatusOK,
		Message: "successfully create a project",
	})
}

func GetProject(c echo.Context) error {
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
		"user_id":    session.UserID,
		"deleted_at": nil,
	}).Decode(&projectPermission); err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "failed to retrieve project-permission",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res.Result{
		Status:  http.StatusOK,
		Message: "successfully getting projects",
		//Details:
	})
}
