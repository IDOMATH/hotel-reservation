package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Booking struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId       primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	RoomId       primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	NumPersons   int                `bson:"numPersons,omitempty" json:"numPersons,omitempty"`
	CheckInDate  time.Time          `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	CheckOutDate time.Time          `bson:"tillDate,omitempty" json:"tillDate,omitempty"`
	Canceled     bool               `bson:"canceled" json:"canceled"`
}
