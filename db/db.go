package db

const (
	DbName = "hotel-reservation"
	TestDbName = "hotel-reservation-test"
	DbUri = "mongodb://localhost:27017"
)

type Store struct {
	User UserStore
	Hotel HotelStore
	Room RoomStore
	Booking BookingStore
}