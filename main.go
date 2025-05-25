package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Body      string             `json:"body" bson:"body"`
	Completed bool               `json:"completed" bson:"completed"`
}

var todos []Todo
var collection *mongo.Collection

func main() {
	fmt.Println("Hello, World!")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	MONGODB_URI := os.Getenv("MONGODB_URI")
	// fmt.Println(MONGODB_URI)
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(MONGODB_URI).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	// Check the connection
	if err != nil {
		log.Fatal("Error connecting to MongoDB", err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	fmt.Println("Connected to MongoDB")

	collection = client.Database("golang_db").Collection("todos")

	// defer client.Disconnect(context.Background())

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	// app.Patch("/api/todos", updateTodo)
	// app.Delete("/api/todos", deleteTodo)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	fmt.Printf("Server starting on port %s\n", PORT)
	log.Fatal(app.Listen(":" + PORT))
}

func getTodos(c *fiber.Ctx) error {
	cursor, err := collection.Find(context.Background(), bson.M{}) // bson.M{} is a filter to find all todos
	if err != nil {
		return c.Status(500).SendString("Error fetching todos")
	}

	defer cursor.Close(context.Background())

	fetchedTodos := make([]Todo, 0) // 初始化为空数组，不是 nil
	if err = cursor.All(context.Background(), &fetchedTodos); err != nil {
		return c.Status(500).SendString("Error decoding todos")
	}

	return c.JSON(fetchedTodos)
}

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)
	if err := c.BodyParser(todo); err != nil {
		return c.Status(400).SendString("Invalid request body")
	}
	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body cannot be empty"})
	}

	fmt.Println(todo)

	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return c.Status(500).SendString("Error creating todo")
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)
}

// func updateTodo(c *fiber.Ctx) error {
// 	return c.SendString("Hello, update todo!")
// }

// func deleteTodo(c *fiber.Ctx) error {
// 	return c.SendString("Hello, delete todo!")
// }
