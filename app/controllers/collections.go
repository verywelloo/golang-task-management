package controllers

import (
	"errors"

	"github.com/patcharp/golib/util"
	d "github.com/verywelloo/3-go-echo-task-management/app/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var DB *mongo.Client = d.ConnectDB()

func GetDatabaseCollection(client *mongo.Client, databaseName string, collectionName string) *mongo.Collection {
	collection := client.Database(databaseName).Collection(collectionName)
	return collection
}

var (
	databaseName = util.GetEnv("DB_NAME", "")

	UserCollection *mongo.Collection = GetDatabaseCollection(DB, databaseName, "users")

	ErrEmptyID   = errors.New("empty ID")
	ErrInvalidID = errors.New("invalid ID")
	ErrNotFound  = errors.New("not found")
)
