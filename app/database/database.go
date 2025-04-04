package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	s "github.com/verywelloo/3-go-echo-task-management/app/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := s.GetEnv("MONGOURI", "")
	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	// ping check database
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("err in ping client %v", err)
		return nil
	}

	parsedURI, err := url.Parse(uri)
	if err != nil {
		log.Fatal("Failed to parse URI:", err)
	}

	log.Fatal("Connected to MongoDB ---> ", fmt.Sprintf("%s:%s, %s", parsedURI.Hostname(), parsedURI.Port(), s.GetEnv("DB_NAME", "")), err)

	return Client
}
