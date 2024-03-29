package main

import (
	"go-chat-app/authentication"
	"go-chat-app/communication"
	"go-chat-app/database"
	"log"
	"os"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env vars
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error while loading env: \n%v\n", err)
	}

	// Connect to postgresql database
	database.DB, err = database.InitDb(database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		SSLMode:  os.Getenv("SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	})

	if err != nil {
		log.Fatalf("Error while opening db: \n%v\n", err)
	}

	database.DB.Table("users").AutoMigrate(&authentication.User{})
	database.DB.Table("sessions").AutoMigrate(&authentication.Session{})
	database.DB.Table("messages").AutoMigrate(&communication.Message{})
	database.DB.Table("chat_rooms").AutoMigrate(&communication.ChatRoom{})

	router := fiber.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:9000",
		AllowMethods:     "GET, POST, PATCH, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Origin, Content-Type",
		AllowCredentials: true,
	}))

	router.Use(LoggerMiddleware)

	api := router.Group("/api")
	// Authentication
	api.Post("/register", authentication.Register)
	api.Post("/login", authentication.Login)
	api.Get("/logout", authentication.Logout)
	api.Get("/code/:code", authentication.EmailCodeVerifier)

	// Rooms
	api.Get("/rooms", communication.GetChatRooms)
	api.Get("/rooms/:id", communication.GetChatRoomByID)
	api.Patch("/rooms", communication.EditChatRoom)
	api.Post("/rooms", communication.CreateChatRoom)

	// Users
	api.Get("/users", communication.GetUsers)
	api.Get("/users/:id", communication.GetUserByID)
	api.Get("/getUserData", communication.GetUserData)

	// Messages
	api.Put("/messages", communication.GetMessages)
	api.Delete("message/:id", communication.DeleteMessage)

	// Websocket endpoints
	api.Get("/socket/:id", websocket.New(communication.SendMessage))

	// Start server
	err = router.Listen(":7000")
	if err != nil {
		log.Fatalf("Error while starting server: \n%v\n", err)
	}
}

func LoggerMiddleware(ctx *fiber.Ctx) error {
	log.Printf("New request To %v with method %v\n", ctx.Path(), ctx.Method())
	return ctx.Next()
}
