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

	var assigneeIDs []primitive.ObjectID
	for _, assignee := range payload.Assignee {
		assigneeID, _ := primitive.ObjectIDFromHex(assignee)

		assigneeIDs = append(assigneeIDs, assigneeID)
	}

	insert := m.Task{
		ID:        primitive.NewObjectID(),
		TaskName:  payload.TaskName,
		ProjectID: projectID,
		Assignee:  assigneeIDs,
		StartDate: startDate,
		EndDate:   endDate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

	projectIDStr := c.Param("project_id")

	projectID, err := primitive.ObjectIDFromHex(projectIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, res.Result{
			Status:  http.StatusBadRequest,
			Message: "invalid project id",
			Details: err.Error(),
		})
	}

	taskFilter := []bson.M{
		{
			"$match": bson.M{
				"project_id": projectID,
			},
		},

		{"$sort": bson.M{"created_at": -1}},

		// look up assignee
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "assignee",
				"foreignField": "_id",
				"as":           "assignee_details",
			},
		},
		{
			"$addFields": bson.M{
				"assignee_details": bson.M{
					"$map": bson.M{
						"input": "$assignee_details",
						"as":    "user",
						"in": bson.M{
							"_id":           "$$user._id",
							"employee_id":   "$$user.employee_id",
							"position_name": "$$user.position_name",
							"nickname":      "$$user.nickname",
							"full_name": bson.M{
								"$concat": []interface{}{
									"$$user.title_th", "",
									"$$user.first_name_th", " ",
									"$$user.last_name_th",
								},
							},
						},
					},
				},
			},
		},
	}

	// taskFilter := bson.M{
	// 	"project_id": projectID,
	// }

	// option := options.Find().SetSort(bson.M{"created_at": -1})

	cur, err := taskCollection.Aggregate(ctx, taskFilter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Result{
			Status:  http.StatusInternalServerError,
			Message: "failed to retrieve tasks",
			Details: err.Error(),
		})
	}

	// var tasks []m.Task
	// if err := cur.All(ctx, &tasks); err != nil {
	// 	return c.JSON(http.StatusInternalServerError, res.Result{
	// 		Status:  http.StatusInternalServerError,
	// 		Message: "failed to decode tasks",
	// 		Details: err.Error(),
	// 	})
	// }

	var response []res.GetTaskResponse
	for _, t := range tasks {
		result := res.GetTaskResponse{
			ID:        t.ID.Hex(),
			TaskName:  t.TaskName,
			ProjectID: t.ProjectID.Hex(),
			StartDate: t.StartDate.Format("2006-01-02"),
			EndDate:   t.EndDate.Format("2006-01-02"),
		}

		response = append(response, result)
	}

	return c.JSON(http.StatusOK, res.Result{
		Status:  http.StatusOK,
		Message: "successfully get tasks",
		Details: response,
	})
}
