package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status int

const (
	NotStarted Status = 1
	InProgress Status = 2
	Completed  Status = 3
)

type Project struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Status      Status             `json:"status" bson:"status"`
	ClientID    primitive.ObjectID `json:"clientId" bson:"clientId"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}
