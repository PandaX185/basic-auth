package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	FullName string             `bson:"fullname"`
	Email    string             `bson:"email"`
	Username string             `bson:"username,omitempty,required"`
	Password string             `bson:"password,omitempty,required"`
}
