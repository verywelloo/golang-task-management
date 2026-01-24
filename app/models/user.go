package models

import "time"

type User struct {
	Email     string    `bson:"email"`
	Name      string    `bson:"name"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
