package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type Url struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Link    string             `bson:"link,omitempty"`
	Parsed  bool               `bson:"parsed,omitempty"`
	Error   string             `bson:"error,omitempty"`
	Content string             `bson:"content,omitempty"`
	Rank    float64            `bson:"rank,omitempty"`
}
