package main

import (
	"context"
	"log"
	"net/http"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Todo struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Completed bool          `json:"completed"`
	Body      string        `json:"body"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	MONGO_URI := os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	collection := client.Database("todo").Collection("todos")

	app := fiber.New()

	app.Get("/list", func(c *fiber.Ctx) error {
		var todoList []Todo
		cursor, err := collection.Find(context.Background(), bson.M{})
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var todo Todo
			if err := cursor.Decode(&todo); err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
			todoList = append(todoList, todo)
		}

		return c.JSON(todoList)
	})

	app.Post("/add", func(c *fiber.Ctx) error {
		var todo Todo
		if err := c.BodyParser(&todo); err != nil {
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}

		if todo.Body == "" {
			return c.Status(http.StatusBadRequest).SendString("Todo body cannot be empty")
		}

		_, err := collection.InsertOne(context.Background(), todo)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		todo.ID = bson.NewObjectID()

		return c.SendStatus(http.StatusCreated)
	})

	app.Put("/update/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var todo Todo
		if err := c.BodyParser(&todo); err != nil {
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}

		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("Invalid ID format")
		}

		update := bson.M{
			"$set": bson.M{
				"completed": todo.Completed,
				"body":      todo.Body,
			},
		}

		result, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		if result.MatchedCount == 0 {
			return c.Status(http.StatusNotFound).SendString("Todo not found")
		}

		return c.SendStatus(http.StatusOK)
	})
	
	app.Delete("/delete/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("Invalid ID format")
		}

		result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		if result.DeletedCount == 0 {
			return c.Status(http.StatusNotFound).SendString("Todo not found")
		}

		return c.SendStatus(http.StatusOK)
	})

	log.Fatal(app.Listen(":" + PORT))

}
