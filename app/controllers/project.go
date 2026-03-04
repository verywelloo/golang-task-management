package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	//req "github.com/verywelloo/3-go-echo-task-management/app/dto/request"
	res "github.com/verywelloo/3-go-echo-task-management/app/dto/response"
	s "github.com/verywelloo/3-go-echo-task-management/app/services"
)

func CreateProject(c echo.Context) error {
	userCollection = s.AppInstance.Collections.Users

	if session, err := s.GetSessionCache(c); err != nil {
		return c.JSON(http.StatusUnauthorized, res.Result{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
			Details: err.Error(),
		})
	}

	var projectPerm 
	if userCollection.FindOne(c.Request().Context(),bson.M{
		"user_id": session.ID,
		"deleted_at": nil,
	})

	return nil
}
