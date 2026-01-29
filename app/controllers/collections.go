package controllers

import (
	"log"

	s "github.com/verywelloo/3-go-echo-task-management/app/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	var err error
	DB, err = s.InitEnvironment()
	if err != nil {
		log.Fatal(err)
	}
}

var (
	databaseName = s.GetEnv("DB_NAME", "")

	DB *mongo.Client

	UserCollection *mongo.Collection = s.GetDatabaseCollection(DB, databaseName, "users")
)

// import (
// 	"errors"

// 	"github.com/patcharp/golib/util"
// 	d "github.com/verywelloo/3-go-echo-task-management/app/database"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// var DB *mongo.Client = d.ConnectDB()

// func GetDatabaseCollection(client *mongo.Client, databaseName string, collectionName string) *mongo.Collection {
// 	collection := client.Database(databaseName).Collection(collectionName)
// 	return collection
// }

// var (
// 	databaseName = util.GetEnv("DB_NAME", "")

// 	UserCollection *mongo.Collection = GetDatabaseCollection(DB, databaseName, "users")

// 	ErrEmptyID   = errors.New("empty ID")
// 	ErrInvalidID = errors.New("invalid ID")
// 	ErrNotFound  = errors.New("not found")
// )
