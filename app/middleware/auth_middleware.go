package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	req "github.com/verywelloo/3-go-echo-task-management/app/dto/request"
	res "github.com/verywelloo/3-go-echo-task-management/app/dto/response"
	m "github.com/verywelloo/3-go-echo-task-management/app/models"
	s "github.com/verywelloo/3-go-echo-task-management/app/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, res.Result{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, res.Result{
				Status:  http.StatusUnauthorized,
				Message: "invalid authorization format",
			})
		}

		// decode
		accessToken := tokenParts[1]
		claims, err := s.DecodeAccessToken(accessToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, res.Result{
				Status:  http.StatusUnauthorized,
				Message: "invalid token",
			})
		}

		// get session from redis
		sessionKey, err := s.SessionKey(claims.ID)
		if err != nil {
			fmt.Printf("\ncannot get session key\n")
			return c.JSON(http.StatusInternalServerError, res.Result{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Details: err.Error(),
			})
		}

		var session req.CacheSession
		if err := s.GetRedis(c, sessionKey, &session); err != nil {
			return c.JSON(http.StatusUnauthorized, res.Result{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
				Details: err.Error(),
			})
		}

		//validate session
		var errorStr string
		var sessionErr bool
		switch {
		case session.UserID != claims.Subject:
			sessionErr = true
			errorStr = "session_id"
		case session.Ip != c.RealIP():
			sessionErr = true
			errorStr = "session_ip"
		case session.Agent != c.Request().UserAgent():
			sessionErr = true
			errorStr = "session_agent"
		}

		if sessionErr {
			fmt.Printf("\n\ndata not match in %v\n\n", errorStr)
			return c.JSON(http.StatusUnauthorized, res.Result{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
		}

		//set authenticated context
		var authKey = m.ContextKey{}
		request := c.Request()
		authCtx := context.WithValue(request.Context(), authKey, claims)
		c.SetRequest(request.WithContext(authCtx))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userCollection := s.AppInstance.Collections.Users

		userID, err := primitive.ObjectIDFromHex(session.UserID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, res.Result{
				Status:  http.StatusInternalServerError,
				Message: "invalid user_id",
				Details: err.Error(),
			})
		}

		var user m.User
		err = userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.JSON(http.StatusNotFound, res.Result{
					Status:  http.StatusNotFound,
					Message: "user not found",
					Details: err.Error(),
				})
			}
			return c.JSON(http.StatusInternalServerError, res.Result{
				Status:  http.StatusInternalServerError,
				Message: "failed to retrieve user",
				Details: err.Error(),
			})
		}

		return next(c)
	}
}
