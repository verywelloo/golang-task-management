package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	ShutdownCtx context.Context
}

type Collections struct {
	Users    *mongo.Collection
	Projects *mongo.Collection
}

type App struct {
	DB          *mongo.Client
	Collections *Collections
}

var AppService *Service
var AppInstance *App

func InitializeData(ctx context.Context) error {
	// set mongo connection
	db, err := InitEnvironment()
	if err != nil {
		fmt.Printf("\n%v\n", err)
		return err
	}

	// init collection
	if err := InitCollection(db, ctx); err != nil {
		fmt.Printf("\n%v\n", err)
		return err
	}

	AppService = &Service{
		ShutdownCtx: ctx,
	}
	AppInstance = &App{
		DB:          db,
		Collections: NewCollections(db),
	}
	return nil
}

func InitEnvironment() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// created client, open connection pool
	client, err := mongo.Connect(
		ctx, options.Client().ApplyURI(GetEnv("MONGOURI", "")), // set config
	)
	if err != nil {
		return nil, err
	}

	// check connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, err // client = mongodb configs
}

func InitCollection(client *mongo.Client, ctx context.Context) error {
	db := client.Database(GetEnv("DB_NAME", ""))

	collections := []string{
		"projects",
		"users",
		"tasks",
	}

	// create collection in mongo
	for _, coll := range collections {
		err := db.CreateCollection(ctx, coll)
		if err != nil {
			// skip duplicate error
			if !mongo.IsDuplicateKeyError(err) && !strings.Contains(err.Error(), "NamespaceExists") {
				fmt.Printf("\n%v\n", err)
				return err
			}
		}
	}

	return nil
}

func GetDatabaseCollection(client *mongo.Client, databaseName string, collectionName string) *mongo.Collection {
	collection := client.Database(databaseName).Collection(collectionName)
	return collection
}

func NewCollections(db *mongo.Client) *Collections {
	database := db.Database(GetEnv("DB_NAME", ""))

	return &Collections{
		Users:    database.Collection("users"),
		Projects: database.Collection("projects"),
	}
}
