package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
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
	Redis       *redis.Client
}

var AppService *Service
var AppInstance *App

func InitializeData(ctx context.Context) error {
	// set mongo connection
	db, err := InitEnvironment()
	if err != nil {
		fmt.Printf("\n%w\n", err)
		return err
	}

	// init collection
	if err := InitCollection(db, ctx); err != nil {
		fmt.Printf("\n%w\n", err)
		return err
	}

	// connect redis
	redis, err := connectRedis(ctx)
	if err != nil {
		fmt.Printf("\n%w\n", err)
	}

	// // keep for whole server
	AppService = &Service{
		ShutdownCtx: ctx,
	}

	// keep for setting later
	AppInstance = &App{
		DB:          db,
		Collections: NewCollections(db),
		Redis:       redis,
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

func connectRedis(ctx context.Context) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s%s", GetEnv("REDIS_HOST", "localhost"), GetEnv("REDIS_PORT", "6379")),
		Password: GetEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	fmt.Printf("\nConnect to Redis ---> %v\n", redisClient.Options().Addr)
	return redisClient, nil
}
