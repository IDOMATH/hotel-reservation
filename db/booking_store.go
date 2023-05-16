package db

import (
	"context"
	"github.com/IDOMATH/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	GetBookingById(context.Context, string) (*types.Booking, error)
	UpdateBooking(context.Context, string, bson.M) error
}

type MongoBookingStore struct {
	client     *mongo.Client
	collection *mongo.Collection

	BookingStore
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client:     client,
		collection: client.Database(DbName).Collection("bookings"),
	}
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id string, update bson.M) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	m := bson.M{"$set": update}
	_, err = s.collection.UpdateByID(ctx, oid, m)
	return err
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := s.collection.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.Id = resp.InsertedID.(primitive.ObjectID)
	return booking, nil
}
