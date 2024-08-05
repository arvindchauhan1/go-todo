package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello World!")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env")
	}

	// PORT := 	os.Getenv("PORT")
	MONGO_URL := os.Getenv("MONGO_URL")

	clientOptions := options.Client().ApplyURI(MONGO_URL)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mongoDB...")

	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	//get all todos
	app.Get("/api/todos", GetTodos)

	//create todos
	app.Post("/api/todos", CreateTodo)

	// update todo
	app.Put("/api/todos/:id", UpdateTodo)

	//delete todo
	app.Delete("/api/todos/:id", DeleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	log.Fatal(app.Listen("0.0.0.0:" + port))

}
