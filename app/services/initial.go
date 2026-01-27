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

var AppService *Service

func InitializeData(ctx context.Context) error {
	// set mongo connection
	db, err := initEnvironment()
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
	return nil
}

func initEnvironment() (*mongo.Client, error) {
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
