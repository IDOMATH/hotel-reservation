package main

import (
	"context"
	"flag"
	"github.com/IDOMATH/hotel-reservation/api"
	"github.com/IDOMATH/hotel-reservation/api/middleware"
	"github.com/IDOMATH/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const dbURI = "mongodb://localhost:27017"
const dbName = "hotel-reservation"
const userCollection = "users"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	portNumber := flag.String("port", ":8080", "The listen port for the server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal(err)
	}

	// Handlers
	var (
		userStore    = db.NewMongoUserStore(client, db.DbName)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}

		authHandler    = api.NewAuthHandler(userStore)
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)

		app   = fiber.New(config)
		auth  = app.Group("/api")
		apiV1 = app.Group("/api/v1")
		admin = apiV1.Group("/admin", middleware.AdminAuth)
	)

	// Auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// User handlers
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Post("/user", userHandler.HandlePostUser)

	// Hotel handlers
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// Room handlers
	apiV1.Get("/room", roomHandler.HandleGetRooms)
	apiV1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// Bookings handlers
	apiV1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	// This should probably be a put
	apiV1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// Admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	app.Listen(*portNumber)
}
