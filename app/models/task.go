package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID        primitive.ObjectID   `bson:"_id"`
	TaskName  string               `bson:"task_name"`
	ProjectID primitive.ObjectID   `bson:"project_id"`
	Assignee  []primitive.ObjectID `bson:"assignee"`
	StartDate time.Time            `bson:"start_date"`
	EndDate   time.Time            `bson:"end_date"`
	CreatedAt time.Time            `bson:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at"`
}
