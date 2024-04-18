package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type Url struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Link    string             `bson:"link"`
	Parsed  bool               `bson:"parsed"`
	Error   string             `bson:"error"`
	Content string             `bson:"content"`
}
