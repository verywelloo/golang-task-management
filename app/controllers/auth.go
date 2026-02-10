package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"unicode"

	"github.com/labstack/echo/v4"
	req "github.com/verywelloo/3-go-echo-task-management/app/dto/request"
	res "github.com/verywelloo/3-go-echo-task-management/app/dto/response"
	m "github.com/verywelloo/3-go-echo-task-management/app/models"
	s "github.com/verywelloo/3-go-echo-task-management/app/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := s.AppInstance.Collections.Users

	var payload req.RegisterPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, res.Result{
			Status:  http.StatusBadRequest,
			Message: "invalid request payload",
			Details: err.Error(),
		})
	}

	var hasUpper, hasLower, hasDigit, lengthCorrect bool
	for _, char := range payload.Password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if len(payload.Password) >= 8 && len(payload.Password) <= 20 {
		lengthCorrect = true
	}

	var errStr string
	if !hasUpper {
		errStr = "password must have a upper case"
	}

	if !hasLower {
		if len(errStr) > 0 {
			errStr = errStr + " , lower case"
		} else {
			errStr = "password must have a lower case"
		}
	}

	if !hasDigit {
		if len(errStr) > 0 {
			errStr = errStr + ", a digit"
		} else {
			errStr = "password must have a digit"
		}
	}

	if !lengthCorrect {
		errStr = errStr + "Length of password must more than 8 and cannot exceed 20"
	}

	// check exists email
	var user m.User
	if err := userCollection.FindOne(ctx, bson.M{"email": payload.Email}).Decode(&user); err == nil {
		return c.JSON(http.StatusBadRequest, res.Result{
			Status:  http.StatusBadRequest,
			Message: "user already exists",
			Details: err.Error(),
		})
	}

	HashPwd, err := s.HashPassword(payload.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "error in hashing password",
			Details: err.Error(),
		})
	}

	if user.Name == "" {
		insert := m.User{
			Email:     payload.Email,
			Name:      payload.Name,
			Password:  HashPwd,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := userCollection.InsertOne(ctx, insert)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, res.Result{
				Status:  http.StatusInternalServerError,
				Message: "error in inserting user",
				Details: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, res.Result{
		Status:  http.StatusOK,
		Message: "error in inserting user",
	})
}

func Login(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := s.AppInstance.Collections.Users

	var payload req.LoginPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, res.Result{
			Status:  http.StatusBadRequest,
			Message: "invalid payload format",
			Details: err.Error(),
		})
	}

	if err := c.Validate(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, res.Result{
			Status:  http.StatusBadRequest,
			Message: "missing require parameter",
			Details: err.Error(),
		})
	}

	var user m.User
	if err := userCollection.FindOne(ctx, bson.M{"email": payload.Email}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusInternalServerError, "user not found")
		}
		return c.JSON(http.StatusInternalServerError, "failed to find a user")
	}

	isPasswordCorrect, err := verifyPassword(payload.Password, user.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, res.Result{
			Status:  http.StatusUnauthorized,
			Message: "error to verify password",
			Details: err.Error(),
		})
	}

	if !isPasswordCorrect {
		return c.JSON(http.StatusUnauthorized, res.Result{
			Status:  http.StatusUnauthorized,
			Message: "password is incorrect",
		})
	}

	sessionID, err := s.GenerateSessionID()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "error to generate session ID",
			Details: err.Error(),
		})
	}

	privateKey, _, err := s.GetRSAKeys(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "cannot get private key",
			Details: err.Error(),
		})
	}

	token, err := s.EncodeAccessToken(sessionID, user.ID.Hex(), user.Name, privateKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "failed to get token",
			Details: err.Error(),
		})
	}

	sessionKey, err := s.SessionKey(sessionID)
	if err != nil {
		fmt.Printf("\ncannot get session key\n")
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
			Details: err.Error(),
		})
	}

	return nil
}

func verifyPassword(candidatePassword, password string) (bool, error) {
	var isPasswordCorrect bool
	hashPwd, err := s.HashPassword(candidatePassword)
	if err != nil {
		return false, err
	}

	if hashPwd == password {
		isPasswordCorrect = true
	}

	return isPasswordCorrect, nil
}
